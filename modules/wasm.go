package modules

import (
	"archive/zip"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/autonomouskoi/akcore"
	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/exe/run"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/akcore/storage/kv"
	svc "github.com/autonomouskoi/akcore/svc/pb"
)

const pluginSuffix = ".akplugin"

func wasmDir(pluginPath string) (fs.FS, error) {
	if !strings.HasSuffix(pluginPath, pluginSuffix) {
		return os.DirFS(pluginPath), nil
	}
	stat, err := os.Stat(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("statting: %w", err)
	}
	fh, err := os.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("opening: %w", err)
	}
	return zip.NewReader(fh, stat.Size())
}

func RegisterWASM(pluginPath string) error {
	plugin, err := wasmDir(pluginPath)
	if err != nil {
		return fmt.Errorf("reading plugin: %w", err)
	}
	b, err := fs.ReadFile(plugin, "manifest.json")
	if err != nil {
		return fmt.Errorf("reading manifest: %w", err)
	}
	manifest := &Manifest{}
	if err := protojson.Unmarshal(b, manifest); err != nil {
		return fmt.Errorf("unmarshalling manifest: %w", err)
	}

	iconBytes, iconType := findIcon(plugin)

	return Register(manifest, &WASM{
		manifest:  manifest,
		basePath:  pluginPath,
		iconBytes: iconBytes,
		iconType:  iconType,
	})
}

type WASM struct {
	modutil.ModuleBase
	manifest  *Manifest
	lock      sync.Mutex
	basePath  string
	bus       *bus.Bus
	kv        kv.KVPrefix
	subs      map[string]chan<- *bus.BusMessage
	in        chan *bus.BusMessage
	wg        sync.WaitGroup
	svc       modutil.Service
	iconBytes []byte
	iconType  string
	http.Handler
}

func (w *WASM) Start(ctx context.Context, deps *modutil.ModuleDeps) error {
	w.Log = deps.Log
	w.bus = deps.Bus
	w.subs = map[string]chan<- *bus.BusMessage{}
	w.kv = deps.KV
	w.in = make(chan *bus.BusMessage, 4)
	w.svc = deps.Svc

	go func() {
		<-ctx.Done()
		w.lock.Lock()
		for topic, c := range w.subs {
			w.bus.Unsubscribe(topic, c)
		}
		w.subs = nil
		w.lock.Unlock()
	}()

	pluginFiles, err := wasmDir(w.basePath)
	if err != nil {
		return fmt.Errorf("reading plugin: %w", err)
	}

	dirEntries, err := fs.ReadDir(pluginFiles, ".")
	if err != nil {
		return fmt.Errorf("reading %s: %w", w.basePath, err)
	}
	wasmFiles := []string{}
	fsPaths := map[string]string{}
	webPath := ""

	for _, de := range dirEntries {
		absPath := de.Name()
		switch {
		case filepath.Ext(de.Name()) == ".wasm":
			wasmFiles = append(wasmFiles, absPath)
		case de.Name() == "data" && de.IsDir():
			fsPaths[absPath] = "/data"
		case de.Name() == "web":
			webPath = absPath
		}
	}

	if len(wasmFiles) == 0 {
		return fmt.Errorf("no wasm files in %s", w.basePath)
	}
	w.Log.Debug("data path", "plugin", w.manifest.Name, "path", deps.StoragePath)

	if webPath != "" {
		sub, err := fs.Sub(pluginFiles, webPath)
		if err != nil {
			return fmt.Errorf("setting web page: %w", err)
		}
		mux := http.NewServeMux()
		mux.Handle("/", http.FileServer(http.FS(sub)))
		if w.manifest.CustomWebDir {
			customWebDir := filepath.Join(deps.StoragePath, "custom-web")
			if err := os.MkdirAll(customWebDir, 0744); err != nil {
				return fmt.Errorf("creating %s: %w", customWebDir, err)
			}
			mux.Handle("/custom-web/", http.StripPrefix("/custom-web/", http.FileServer(http.Dir(customWebDir))))
		}
		w.Handler = mux
	}

	manifest := extism.Manifest{
		AllowedPaths: fsPaths,
	}
	for _, wasmFile := range wasmFiles {
		b, err := fs.ReadFile(pluginFiles, wasmFile)
		if err != nil {
			return fmt.Errorf("reading wasm file %s: %w", wasmFile, err)
		}
		manifest.Wasm = append(manifest.Wasm, extism.WasmData{
			Data: b,
			Name: wasmFile,
		})
	}
	config := extism.PluginConfig{EnableWasi: true}

	plugin, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{
		w.sendFn(),
		w.sendReplyFn(),
		w.waitForReplyFn(),
	})
	if err != nil {
		return fmt.Errorf("initializing plugin: %w", err)
	}
	defer plugin.Close(ctx)

	if _, _, err := plugin.Call("start", nil); err != nil {
		return fmt.Errorf("calling start: %w", err)
	}

	go func() {
		w.wg.Wait()
		close(w.in)
	}()

	for msg := range w.in {
		b, err := proto.Marshal(msg)
		if err != nil {
			return err
		}
		_, replyB, err := plugin.Call("recv", b)
		if err != nil {
			return err
		}
		if len(replyB) == 0 {
			continue
		}
		reply := &bus.BusMessage{}
		if err := proto.Unmarshal(replyB, reply); err != nil {
			return err
		}
		deps.Bus.Send(reply)
	}

	return nil
}

func findIcon(plugin fs.FS) ([]byte, string) {
	dirEntries, err := fs.ReadDir(plugin, ".")
	if err != nil {
		return nil, ""
	}
	for _, de := range dirEntries {
		name := de.Name()
		mimeType := iconType(name)
		if mimeType == "" {
			continue
		}
		b, err := fs.ReadFile(plugin, name)
		if err != nil {
			continue
		}
		return b, mimeType
	}
	return nil, ""
}

func iconType(name string) string {
	if !strings.HasPrefix(name, "icon.") {
		return ""
	}
	switch filepath.Ext(name) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".svg":
		return "image/svg+xml"
	}
	return ""
}

func (w *WASM) handleExternalFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	switch msg.GetType() {
	case int32(bus.ExternalMessageType_HAS_TOPIC_REQ):
		return w.handleHasTopicFromPlugin(msg)
	case int32(bus.ExternalMessageType_SUBSCRIBE_REQ):
		return w.handleSubscribeFromPlugin(msg)
	case int32(bus.ExternalMessageType_UNSUBSCRIBE_REQ):
		return w.handleUnsubscribeFromPlugin(msg)
	case int32(bus.ExternalMessageType_KV_SET_REQ):
		return w.handleKVSetFromPlugin(msg)
	case int32(bus.ExternalMessageType_KV_GET_REQ):
		return w.handleKVGetFromPlugin(msg)
	case int32(bus.ExternalMessageType_KV_LIST_REQ):
		return w.handleKVListFromPlugin(msg)
	case int32(bus.ExternalMessageType_KV_DELETE_REQ):
		return w.handleKVDeleteFromPlugin(msg)
	case int32(bus.ExternalMessageType_LOG_SEND_REQ):
		return w.handleLogSendFromPlugin(msg)
	case int32(svc.MessageType_WEBCLIENT_STATIC_DOWNLOAD_REQ):
		return w.svc.Handle(msg)
	}
	return &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
		Error: &bus.Error{
			Code: int32(bus.CommonErrorCode_INVALID_TYPE),
		},
	}
}

func (w *WASM) handleHasTopicFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	return handle(w, msg, &bus.HasTopicRequest{},
		func(req *bus.HasTopicRequest) (*bus.HasTopicResponse, *bus.Error) {
			return &bus.HasTopicResponse{
				Topic:    req.GetTopic(),
				HasTopic: w.bus.HasTopic(req.GetTopic()),
			}, nil
		})
}

func (w *WASM) handleSubscribeFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	return handle(w, msg, &bus.SubscribeRequest{},
		func(req *bus.SubscribeRequest) (*bus.SubscribeResponse, *bus.Error) {
			w.subscribe(req.GetTopic())
			return &bus.SubscribeResponse{}, nil
		})
}

func (w *WASM) handleUnsubscribeFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	return handle(w, msg, &bus.UnsubscribeRequest{},
		func(req *bus.UnsubscribeRequest) (*bus.UnsubscribeResponse, *bus.Error) {
			w.unsubscribe(req.GetTopic())
			return &bus.UnsubscribeResponse{}, nil
		})
}

func (w *WASM) subscribe(topic string) {
	if topic == "" {
		return
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.subs == nil {
		return // stopping
	}
	if _, present := w.subs[topic]; present {
		return
	}
	w.wg.Add(1)
	in := make(chan *bus.BusMessage)
	w.bus.Subscribe(topic, in)
	w.subs[topic] = in
	go func() {
		for msg := range in {
			select {
			case w.in <- msg:
			default:
			}
		}
		w.wg.Done()
	}()
}

func (w *WASM) unsubscribe(topic string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	in, present := w.subs[topic]
	if present {
		w.bus.Unsubscribe(topic, in)
		delete(w.subs, topic)
	}
}

func (w *WASM) unmarshalAnd(stackPos uint64, p *extism.CurrentPlugin, fn func(*bus.BusMessage)) {
	busMessage, err := p.ReadBytes(stackPos)
	if err != nil {
		w.Log.Error("reading bytes", "error", err.Error())
		return
	}
	msg := &bus.BusMessage{}
	if err := proto.Unmarshal(busMessage, msg); err != nil {
		w.Log.Error("unmarshalling guest message", "error", err.Error())
		return
	}
	fn(msg)
}

func (w *WASM) sendFn() extism.HostFunction {
	return extism.NewHostFunctionWithStack(
		"send",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			w.unmarshalAnd(stack[0], p, func(msg *bus.BusMessage) {
				msg.FromMod = w.manifest.Id
				if msg.GetTopic() == "" {
					w.handleExternalFromPlugin(msg)
				} else {
					w.bus.Send(msg)
				}
			})
			/*
				var err error
				stack[0], err = p.WriteBytes(nil)
				if err != nil {
					w.Log.Error("writing nil bytes", "function", "send", "error", err.Error())
				}
			*/
		},
		[]api.ValueType{api.ValueTypeI64},
		[]api.ValueType{},
	)
}

func (w *WASM) sendReplyFn() extism.HostFunction {
	return extism.NewHostFunctionWithStack(
		"send_reply",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			w.unmarshalAnd(stack[0], p, func(msg *bus.BusMessage) {
				w.bus.SendReply(msg, msg)
			})
		},
		[]api.ValueType{api.ValueTypeI64},
		[]api.ValueType{},
	)
}

func (w *WASM) waitForReplyFn() extism.HostFunction {
	return extism.NewHostFunctionWithStack(
		"wait_for_reply",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			w.unmarshalAnd(stack[0], p, func(msg *bus.BusMessage) {
				timeoutMS := stack[1]
				var reply *bus.BusMessage
				if msg.GetTopic() == "" {
					reply = w.handleExternalFromPlugin(msg)
				} else {
					ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(timeoutMS))
					reply = w.bus.WaitForReply(ctx, msg)
					cancel()
				}
				var b []byte
				var err error
				if reply != nil {
					b, err = proto.Marshal(reply)
					if err != nil {
						w.Log.Error("marshalling guest reply", "error", err.Error())
						return
					}
				}
				stack[0], err = p.WriteBytes(b)
				if err != nil {
					w.Log.Error("writing return bytes", "error", err.Error())
				}
			})
		},
		[]api.ValueType{api.ValueTypeI64, api.ValueTypeI64},
		[]api.ValueType{api.ValueTypeI64},
	)
}

func (w *WASM) handleKVSetFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	return handle(w, msg, &bus.KVSetRequest{},
		func(req *bus.KVSetRequest) (*bus.KVSetResponse, *bus.Error) {
			if err := w.kv.Set(req.GetKey(), req.GetValue()); err != nil {
				return nil, &bus.Error{
					Detail: proto.String(err.Error()),
				}
			}
			return &bus.KVSetResponse{}, nil
		},
	)
}

func (w *WASM) handleKVGetFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	return handle(w, msg, &bus.KVGetRequest{},
		func(req *bus.KVGetRequest) (*bus.KVGetResponse, *bus.Error) {
			v, err := w.kv.Get(req.GetKey())
			if err != nil {
				if errors.Is(err, akcore.ErrNotFound) {
					return nil, &bus.Error{
						Code: int32(bus.CommonErrorCode_NOT_FOUND),
					}
				}
				return nil, &bus.Error{
					Detail: proto.String(err.Error()),
				}
			}
			return &bus.KVGetResponse{
				Key:   req.GetKey(),
				Value: v,
			}, nil
		})
}

func (w *WASM) handleKVListFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	return handle(w, msg, &bus.KVListRequest{},
		func(req *bus.KVListRequest) (*bus.KVListResponse, *bus.Error) {
			matches, err := w.kv.List(req.GetPrefix())
			if err != nil {
				return nil, &bus.Error{
					Detail: proto.String(err.Error()),
				}
			}
			resp := &bus.KVListResponse{
				TotalMatches: uint32(len(matches)),
				Prefix:       req.GetPrefix(),
				Offset:       req.GetOffset(),
				Limit:        req.GetLimit(),
			}
			offset := int(req.GetOffset())
			limit := int(req.GetLimit())
			lenMatches := len(matches)
			if offset >= lenMatches {
				return resp, nil
			}
			if limit == 0 {
				limit = 2<<32 - 1
			}
			for i := 0; i < limit && i < lenMatches; i++ {
				resp.Keys = append(resp.Keys, matches[i+offset])
			}
			return resp, nil
		})
}

func (w *WASM) handleKVDeleteFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	return handle(w, msg, &bus.KVDeleteRequest{},
		func(req *bus.KVDeleteRequest) (*bus.KVDeleteResponse, *bus.Error) {
			if err := w.kv.Delete(req.GetKey()); err != nil {
				return nil, &bus.Error{
					Detail: proto.String(err.Error()),
				}
			}
			return &bus.KVDeleteResponse{}, nil
		})
}

func (w *WASM) handleLogSendFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	return handle(w, msg, &bus.LogSendRequest{},
		func(req *bus.LogSendRequest) (*bus.LogSendResponse, *bus.Error) {
			attrs := make([]slog.Attr, len(req.Args))
			for i, arg := range req.Args {
				attr := slog.Attr{Key: arg.GetKey()}
				switch x := arg.Value.(type) {
				case *bus.LogSendRequest_Arg_String_:
					attr.Value = slog.StringValue(x.String_)
				case *bus.LogSendRequest_Arg_Bool:
					attr.Value = slog.BoolValue(x.Bool)
				case *bus.LogSendRequest_Arg_Int64:
					attr.Value = slog.Int64Value(x.Int64)
				case *bus.LogSendRequest_Arg_Double:
					attr.Value = slog.Float64Value(x.Double)
				case nil:
					attr.Value = slog.AnyValue(nil)
				default:
					attr.Value = slog.StringValue("<unknown type>")
				}
				attrs[i] = attr
			}
			var level slog.Level
			switch req.Level {
			case bus.LogLevel_DEBUG:
				level = slog.LevelDebug
			case bus.LogLevel_INFO:
				level = slog.LevelInfo
			case bus.LogLevel_WARN:
				level = slog.LevelWarn
			case bus.LogLevel_ERROR:
				level = slog.LevelError
			}
			w.Log.LogAttrs(context.Background(), level, req.GetMessage(), attrs...)
			return &bus.LogSendResponse{}, nil
		})
}

func handle[REQ, RESP protoreflect.ProtoMessage](w *WASM, msg *bus.BusMessage, req REQ, fn func(REQ) (RESP, *bus.Error)) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	if reply.Error = w.UnmarshalMessage(msg, req); reply.Error != nil {
		return reply
	}
	var resp protoreflect.ProtoMessage
	resp, reply.Error = fn(req)
	if reply.Error != nil {
		return reply
	}
	w.MarshalMessage(reply, resp)
	return reply
}

func (controller *controller) initWASM() {
	appPath, err := run.AppPath()
	if err != nil {
		return
	}

	installPath := filepath.Join(appPath, "plugins", "install")
	if _, err := os.Stat(installPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			controller.Log.Error("checking plugin install dir", "error", err.Error())
		} else if err := os.MkdirAll(installPath, 0700); err != nil {
			controller.Log.Error("creating plugin install dir", "path", installPath, "error", err.Error())
		}
	}
	dirEnts, err := os.ReadDir(installPath)
	if err != nil {
		controller.Log.Error("reading plugin install dir", "path", installPath, "error", err.Error())
	}
	for _, dirEnt := range dirEnts {
		if !strings.HasSuffix(dirEnt.Name(), pluginSuffix) {
			continue
		}
		zipPath := filepath.Join(installPath, dirEnt.Name())
		if err := RegisterWASM(zipPath); err != nil {
			controller.Log.Error("registering WASM plugin", "path", zipPath, "error", err.Error())
			continue
		}
	}

	devWASMPath := filepath.Join(appPath, "devwasm")
	infh, err := os.Open(devWASMPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			controller.Log.Error("reading devwasm", "path", devWASMPath, "error", err.Error())
		}
		return
	}
	defer infh.Close()
	//controller.Log.Debug("using devwasm", "path", devWASMPath)
	scanner := bufio.NewScanner(infh)
	for scanner.Scan() {
		path := strings.TrimSpace(scanner.Text())
		if path == "" || strings.HasPrefix(path, "#") {
			continue
		}
		if err := RegisterWASM(path); err != nil {
			controller.Log.Error("registering devwasm", "manifest_path", path, "error", err.Error())
			continue
		}
		controller.Log.Debug("loaded devwasm", "path", path)
	}
	if err := scanner.Err(); err != nil {
		controller.Log.Error("scanning devwasm", "path", devWASMPath, "error", err.Error())
		return
	}
}

// Icon returns the icon data and its MIME type. If the plugin includes an icon
// that is returned. If not, a default one is returned.
func (w *WASM) Icon() ([]byte, string, error) {
	if len(w.iconBytes) == 0 {
		return w.ModuleBase.Icon()
	}
	return w.iconBytes, w.iconType, nil
}
