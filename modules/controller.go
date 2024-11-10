package modules

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"path"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/akcore/storage/kv"
)

var (
	ErrDuplicate = errors.New("duplicate module")

	ctrl = &controller{
		modules: map[string]*module{},
		webHandlers: &handler{
			handlers: map[string]http.Handler{},
			mux:      &http.ServeMux{},
		},
	}
)

type controller struct {
	modutil.ModuleBase
	eg          errgroup.Group
	lock        sync.Mutex
	bus         *bus.Bus
	kv          kv.KV
	runningCtx  context.Context
	internalKV  *kv.KVPrefix
	modules     map[string]*module
	webHandlers *handler
	cachePath   string
}

func Register(manifest *Manifest, mod modutil.Module) error {
	return ctrl.Register(manifest, mod)
}

func (controller *controller) Register(manifest *Manifest, mod modutil.Module) error {
	controller.lock.Lock()
	defer controller.lock.Unlock()
	id := manifest.Id
	if _, present := controller.modules[id]; present {
		return ErrDuplicate
	}
	idBytes, err := hex.DecodeString(id)
	if err != nil {
		return fmt.Errorf("decoding ID: %w", err)
	}
	var kvPrefix [8]byte
	copy(kvPrefix[:], idBytes)
	if kvPrefix == [8]byte{} {
		return fmt.Errorf("zero ID: %w", err)
	}
	controller.modules[id] = &module{
		manifest: manifest,
		module:   mod,
		kvPrefix: kvPrefix,
		config:   &Config{},
	}
	return nil
}

func Start(ctx context.Context, deps *modutil.Deps) error {
	return ctrl.Start(ctx, deps)
}

func (controller *controller) Start(ctx context.Context, deps *modutil.Deps) error {
	controller.bus = deps.Bus
	controller.eg = errgroup.Group{}
	controller.Log = deps.Log.With("module", "modules")
	controller.kv = deps.KV
	controller.internalKV = controller.kv.WithPrefix([8]byte{})
	controller.cachePath = deps.CachePath
	controller.runningCtx = ctx

	defer func() {
		// save module configs
		for id, mod := range controller.modules {
			b, err := proto.Marshal(mod.config)
			if err != nil {
				controller.Log.Error("marshalling config", "module_id", id, "error", err.Error())
				continue
			}
			key := []byte("config/" + id)
			haveB, err := controller.internalKV.Get(key)
			if err == nil && bytes.Equal(b, haveB) {
				continue
			}
			if err := controller.internalKV.Set(key, b); err != nil {
				controller.Log.Error("writing config", "module_id", id, "error", err.Error())
			}
		}
	}()

	deps.Web.Handle("/m/", controller.webHandlers)

	controller.eg.Go(func() error {
		in := make(chan *bus.BusMessage, 32)
		deps.Bus.Subscribe(BusTopics_CONTROL.String(), in)
		for msg := range in {
			switch msg.Type {
			case int32(MessageType_TYPE_CHANGE_MODULE_AUTOSTART):
				controller.handleChangeModuleAutostart(msg)
			case int32(MessageType_TYPE_CHANGE_STATE):
				controller.handleChangeState(ctx, msg)
			case int32(MessageType_TYPE_GET_CURRENT_STATES):
				controller.handleGetCurrentStates()
			case int32(MessageType_TYPE_GET_MANIFEST_REQ):
				controller.handleGetManifest(msg)
			default:
				controller.Log.Info("unhandled control message", "type", msg.Type)
			}
		}
		return nil
	})
	controller.eg.Go(func() error { return controller.handleRequests(ctx) })
	controller.eg.Go(func() error { return controller.handleCommand(ctx) })

	if err := deps.Bus.WaitForTopic(ctx, BusTopics_CONTROL.String(), time.Millisecond*10); err != nil {
		return fmt.Errorf("waiting for control topic: %w", err)
	}
	if err := deps.Bus.WaitForTopic(ctx, BusTopics_MODULE_COMMAND.String(), time.Millisecond*10); err != nil {
		return fmt.Errorf("waiting for command topic: %w", err)
	}
	if err := controller.initModules(ctx); err != nil {
		return err
	}

	return controller.eg.Wait()
}

func (controller *controller) handleChangeModuleAutostart(msg *bus.BusMessage) {
	cma := &ChangeModuleAutostart{}
	if err := proto.Unmarshal(msg.GetMessage(), cma); err != nil {
		controller.Log.Error("unmarshalling ChangeModuleAutostart", "error", err.Error())
		return
	}
	module, ok := controller.modules[cma.ModuleId]
	if !ok {
		return
	}
	module.config.AutomaticStart = cma.Autostart
	module.sendState()
}

func (controller *controller) handleChangeState(ctx context.Context, msg *bus.BusMessage) {
	cs := &ChangeModuleState{}
	if err := proto.Unmarshal(msg.GetMessage(), cs); err != nil {
		controller.Log.Error("unmarshalling ChangeModuleState", "error", err.Error())
		return
	}
	switch cs.ModuleState {
	case ModuleState_STARTED:
		controller.startModule(ctx, cs.ModuleId)
	case ModuleState_STOPPED:
		controller.stopModule(cs.ModuleId)
	default:
		controller.Log.Error("unhandled module state",
			"module_id", cs.ModuleId, "state", cs.ModuleState)
	}
}

func (controller *controller) startModule(ctx context.Context, id string) {
	mod, present := controller.modules[id]
	if !present {
		controller.Log.Error("starting invalid module", "id", id)
		return
	}
	if gotLock := mod.lock.TryLock(); !gotLock {
		controller.Log.Error("starting already running module", "id", id)
		return
	}
	mod.lock.Unlock()

	controller.eg.Go(func() error {
		ctx, mod.cancel = context.WithCancel(ctx)
		defer mod.cancel()
		mod.lock.Lock()
		defer mod.lock.Unlock()
		mod.setState(ModuleState_STARTED)
		controller.Log.Info("starting", "module", id)
		defer controller.Log.Debug("exiting", "module", id)
		if handler, ok := mod.module.(http.Handler); ok {
			path := path.Join("/m", mod.manifest.Name) + "/"
			mod.deps.Log.Debug("registering web handler", "path", path)
			controller.webHandlers.Handle(path, handler)
			defer func() {
				mod.deps.Log.Debug("deregistering web handler", "path", path)
				controller.webHandlers.Remove(path)
			}()
		}
		err := mod.module.Start(ctx, mod.deps)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				mod.setState(ModuleState_STOPPED)
				return nil
			}
			mod.deps.Log.Error("starting", "error", err.Error())
			mod.setState(ModuleState_FAILED)
			return err
		}
		mod.setState(ModuleState_FINISHED)
		return nil
	})
}

func (controller *controller) stopModule(id string) {
	mod, present := controller.modules[id]
	if !present {
		controller.Log.Error("can't stop unregistered module", "id", id)
		return
	}
	mod.cancel()
}

func (controller *controller) handleGetCurrentStates() {
	for _, mod := range controller.modules {
		mod.sendState()
	}
}

func (controller *controller) handleGetManifest(msg *bus.BusMessage) {
	gmReq := &GetManifestRequest{}
	if err := proto.Unmarshal(msg.GetMessage(), gmReq); err != nil {
		controller.Log.Error("unmarshalling GetManifestRequest", "error", err.Error())
		return
	}
	mod, present := controller.modules[gmReq.ModuleId]
	if !present {
		busErr := bus.Error{
			Code:   int32(bus.CommonErrorCode_NOT_FOUND),
			Detail: &gmReq.ModuleId,
		}
		resp := &bus.BusMessage{
			Error: &busErr,
		}
		controller.bus.SendReply(msg, resp)
		return
	}
	b, err := proto.Marshal(&GetManifestResponse{
		Manifest: mod.manifest,
	})
	if err != nil {
		controller.Log.Error("marshalling GetManifestResponse", "error", err.Error())
		return
	}
	resp := &bus.BusMessage{
		Type:    int32(MessageType_TYPE_GET_MANIFEST_RESP),
		Message: b,
	}
	controller.bus.SendReply(msg, resp)
}
