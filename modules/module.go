package modules

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore"
	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
)

// module holds the details for a registered module and manages its lifecycle
type module struct {
	manifest    *Manifest
	state       ModuleState
	stateDetail string
	deps        *modutil.ModuleDeps
	module      modutil.Module
	lock        sync.Mutex
	cancel      func()
	kvPrefix    [8]byte
	config      *Config
}

// set the modules state and send the new state as an event
func (m *module) setState(newState ModuleState) error {
	if m.state == newState {
		return nil
	}
	m.state = newState
	return m.sendState()
}

// send the module's current state as an event so the UI and other modules can
// react as desired
func (m *module) sendState() error {
	msg := &bus.BusMessage{
		Topic: BusTopics_MODULE_EVENT.String(),
		Type:  int32(MessageTypeEvent_MODULE_CURRENT_STATE),
	}
	var err error
	msg.Message, err = proto.Marshal(&ModuleCurrentStateEvent{
		ModuleId:    m.manifest.Id,
		ModuleState: m.state,
		Config:      m.config,
		StateDetail: m.stateDetail,
	})
	if err != nil {
		return err
	}
	m.deps.Bus.Send(msg)
	return nil
}

// perform initialization for registered modules on the controller
func (controller *controller) initModules(ctx context.Context) error {
	for id, mod := range controller.modules {
		if ctx.Err() != nil {
			return nil
		}
		// create the module's dependencies. Not all deps will be used by all
		// modules
		mod.deps = &modutil.ModuleDeps{
			Bus:         controller.bus,
			KV:          *controller.kv.WithPrefix(mod.kvPrefix),
			Log:         controller.Log.With("module", id),
			CachePath:   filepath.Join(controller.cachePath, "AutonomousKoi", id),
			StoragePath: filepath.Join(controller.storagePath, mod.manifest.Name),
		}

		// set the state as unstarted
		mod.setState(ModuleState_UNSTARTED)
		// retrieve the config for this specific module (e.g. autostart)
		configB, err := controller.internalKV.Get([]byte("config/" + id))
		if err == nil {
			if err = proto.Unmarshal(configB, mod.config); err != nil {
				return fmt.Errorf("unmarshalling config for %s: %w", id, err)
			}
		} else if !errors.Is(err, akcore.ErrNotFound) {
			return fmt.Errorf("getting config for %s: %w", id, err)
		}
		if !mod.config.AutomaticStart {
			continue // don't start automatically, proceed to the next
		}

		// start the module by sending the message to start the module
		msg := &bus.BusMessage{
			Topic: BusTopics_MODULE_COMMAND.String(),
			Type:  int32(MessageTypeCommand_MODULE_STATE_SET_REQ),
		}
		controller.MarshalMessage(msg, &ModuleStateSetRequest{
			ModuleId: id,
			State:    ModuleState_STARTED,
		})
		if msg.Error != nil {
			continue
		}
		controller.bus.Send(msg)
	}
	return nil
}
