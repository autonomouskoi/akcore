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
	// ErrDuplicate indicates a duplicate registration for a module
	ErrDuplicate = errors.New("duplicate module")

	ctrl = &controller{
		modules: map[string]*module{},
		webHandlers: &handler{
			handlers: map[string]http.Handler{},
			mux:      &http.ServeMux{},
		},
	}
)

// controller is the backbone of module handling. Modules are registered with
// a controller. The controller manages the lifecycle of modules.
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
	storagePath string
}

// Register a module with the default controller
func Register(manifest *Manifest, mod modutil.Module) error {
	return ctrl.Register(manifest, mod)
}

// Register a module with the controller. Registered modules can be started and
// stopped via requests on the bus. If the manifest contains the ID of a module
// that's been registered already ErrDuplicate is returned.
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

// Start the default controller.
func Start(ctx context.Context, deps *modutil.Deps) error {
	return ctrl.Start(ctx, deps)
}

// Start this controller. Registered modules with autostart true will be started
// automatically. Returns when ctx is cancelled.
func (controller *controller) Start(ctx context.Context, deps *modutil.Deps) error {
	controller.bus = deps.Bus
	controller.eg = errgroup.Group{}
	controller.Log = deps.Log.With("module", "modules")
	controller.kv = deps.KV
	controller.internalKV = controller.kv.WithPrefix([8]byte{})
	controller.cachePath = deps.CachePath
	controller.storagePath = deps.StoragePath
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

	controller.eg.Go(func() error { return controller.handleRequests(ctx) })
	controller.eg.Go(func() error { return controller.handleCommand(ctx) })
	if err := deps.Bus.WaitForTopic(ctx, BusTopics_MODULE_COMMAND.String(), time.Millisecond*10); err != nil {
		return fmt.Errorf("waiting for command topic: %w", err)
	}
	if err := controller.initModules(ctx); err != nil {
		return err
	}

	return controller.eg.Wait()
}

// attempt to start a registered module, referred to by its ID.
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
		// create a context specifically to stop this module as needed
		ctx, mod.cancel = context.WithCancel(ctx)
		defer mod.cancel()
		mod.lock.Lock()
		defer mod.lock.Unlock()
		mod.setState(ModuleState_STARTED)
		controller.Log.Info("starting", "module", id)
		defer controller.Log.Debug("exiting", "module", id)

		// handle webhook calls for this module
		wwh := webhooksHandler(mod.deps.Bus, mod.manifest.Id)
		wwhPath := path.Join("/m", mod.manifest.Id, "_webhook")
		controller.webHandlers.Handle(wwhPath, wwh)
		defer controller.webHandlers.Remove(wwhPath)

		// if the module is an http.Handler, serve its content from a path
		// based on the module's name
		if handler, ok := mod.module.(http.Handler); ok {
			partialPath := path.Join("/m", mod.manifest.Name)
			path := partialPath + "/"
			mod.deps.Log.Debug("registering web handler", "path", path)
			controller.webHandlers.Handle(path, http.StripPrefix(partialPath, handler))
			defer func() {
				mod.deps.Log.Debug("deregistering web handler", "path", path)
				controller.webHandlers.Remove(path)
			}()
		}
		// start the module and block
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

// stop a module by cancelling its context
func (controller *controller) stopModule(id string) {
	mod, present := controller.modules[id]
	if !present {
		controller.Log.Error("can't stop unregistered module", "id", id)
		return
	}
	mod.cancel()
}
