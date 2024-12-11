package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/datastruct/slices"
	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
)

type WASMManifest struct {
	Name        string
	ID          string
	Description string
	WASMFiles   []string
}

func RegisterWASM(manifestPath string) error {
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("reading manifest: %w", err)
	}
	manifest := WASMManifest{}
	if err := json.Unmarshal(b, &manifest); err != nil {
		return fmt.Errorf("unmarshalling manifest: %w", err)
	}
	dir := filepath.Dir(manifestPath)
	return Register(&Manifest{
		Id:          manifest.ID,
		Name:        manifest.Name,
		Description: manifest.Description,
	}, &wasm{
		wasmFiles: slices.Map(manifest.WASMFiles, func(relPath string) string {
			return filepath.Join(dir, relPath)
		}),
	})
}

type wasm struct {
	modutil.ModuleBase
	lock      sync.Mutex
	wasmFiles []string
	bus       *bus.Bus
	subs      map[string]chan<- *bus.BusMessage
	in        chan *bus.BusMessage
	wg        sync.WaitGroup
}

func (w *wasm) Start(ctx context.Context, deps *modutil.ModuleDeps) error {
	w.Log = deps.Log
	w.bus = deps.Bus
	w.subs = map[string]chan<- *bus.BusMessage{}
	w.in = make(chan *bus.BusMessage, 4)
	w.wg.Add(1)
	w.subscribe("fake_topic")

	go func() {
		<-ctx.Done()
		w.lock.Lock()
		for topic, c := range w.subs {
			w.bus.Unsubscribe(topic, c)
		}
		w.subs = nil
		w.lock.Unlock()
	}()
	go func() {
		w.wg.Wait()
		close(w.in)
	}()

	manifest := extism.Manifest{
		Wasm: slices.Map(w.wasmFiles, func(path string) extism.Wasm {
			return extism.WasmFile{Path: path}
		}),
	}
	config := extism.PluginConfig{EnableWasi: true}

	plugin, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{
		w.waitForReplyFn(),
	})
	if err != nil {
		return fmt.Errorf("initializing plugin: %w", err)
	}
	defer plugin.Close()

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

func (w *wasm) handleExternalFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	switch msg.GetType() {
	case int32(bus.ExternalMessageType_SUBSCRIBE_REQ):
		return w.handleSubscribeFromPlugin(msg)
	case int32(bus.ExternalMessageType_UNSUBSCRIBE_REQ):
		return w.handleUnsubscribeFromPlugin(msg)
	}
	return &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
		Error: &bus.Error{
			Code: int32(bus.CommonErrorCode_INVALID_TYPE),
		},
	}
}

func (w *wasm) handleSubscribeFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	sr := &bus.SubscribeRequest{}
	if reply.Error = w.UnmarshalMessage(msg, sr); reply.Error != nil {
		return reply
	}
	w.subscribe(sr.GetTopic())
	w.MarshalMessage(reply, &bus.SubscribeResponse{})
	return reply
}

func (w *wasm) handleUnsubscribeFromPlugin(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	ur := &bus.UnsubscribeRequest{}
	if reply.Error = w.UnmarshalMessage(msg, ur); reply.Error != nil {
		return reply
	}
	w.unsubscribe(ur.GetTopic())
	w.MarshalMessage(reply, &bus.UnsubscribeResponse{})
	return reply
}

func (w *wasm) subscribe(topic string) {
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
			w.in <- msg
		}
		w.wg.Done()
	}()
}

func (w *wasm) unsubscribe(topic string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	in, present := w.subs[topic]
	if present {
		w.bus.Unsubscribe(topic, in)
		delete(w.subs, topic)
	}
}

func (w *wasm) waitForReplyFn() extism.HostFunction {
	return extism.NewHostFunctionWithStack(
		"wait_for_reply",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			busMessage, err := p.ReadBytes(stack[0])
			if err != nil {
				w.Log.Error("reading bytes", "error", err.Error())
				return
			}
			timeoutMS := stack[1]
			msg := &bus.BusMessage{}
			if err := proto.Unmarshal(busMessage, msg); err != nil {
				w.Log.Error("unmarshalling guest message", "error", err.Error())
				return
			}
			var reply *bus.BusMessage
			if msg.GetTopic() == "" {
				reply = w.handleExternalFromPlugin(msg)
			} else {
				ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(timeoutMS))
				reply = w.bus.WaitForReply(ctx, msg)
				cancel()
			}
			var b []byte
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
		},
		[]api.ValueType{api.ValueTypeI64, api.ValueTypeI64},
		[]api.ValueType{api.ValueTypeI64},
	)
}
