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

func (m *module) setState(newState ModuleState) error {
	if m.state == newState {
		return nil
	}
	m.state = newState
	return m.sendState()
}

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

func (controller *controller) initModules(ctx context.Context) error {
	for id, mod := range controller.modules {
		if ctx.Err() != nil {
			return nil
		}
		mod.deps = &modutil.ModuleDeps{
			Bus:         controller.bus,
			KV:          *controller.kv.WithPrefix(mod.kvPrefix),
			Log:         controller.Log.With("module", id),
			CachePath:   filepath.Join(controller.cachePath, "AutonomousKoi", id),
			StoragePath: filepath.Join(controller.storagePath, mod.manifest.Name),
		}

		mod.setState(ModuleState_UNSTARTED)
		configB, err := controller.internalKV.Get([]byte("config/" + id))
		if err == nil {
			if err = proto.Unmarshal(configB, mod.config); err != nil {
				return fmt.Errorf("unmarshalling config for %s: %w", id, err)
			}
		} else if !errors.Is(err, akcore.ErrNotFound) {
			return fmt.Errorf("getting config for %s: %w", id, err)
		}
		if !mod.config.AutomaticStart {
			continue
		}

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
