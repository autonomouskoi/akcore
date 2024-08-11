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
	manifest *Manifest
	state    ModuleState
	deps     *modutil.ModuleDeps
	module   modutil.Module
	lock     sync.Mutex
	cancel   func()
	kvPrefix [8]byte
	config   *Config
}

func (m *module) setState(newState ModuleState) error {
	if m.state == newState {
		return nil
	}
	m.state = newState
	return m.sendState()
}

func (m *module) sendState() error {
	cms, err := proto.Marshal(&CurrentModuleState{
		ModuleId:    m.manifest.Id,
		ModuleState: m.state,
		Config:      m.config,
	})
	if err != nil {
		return fmt.Errorf("marshalling CurrentModuleState: %w", err)
	}
	m.deps.Bus.Send(&bus.BusMessage{
		Topic:   BusTopics_STATE.String(),
		Type:    int32(MessageType_TYPE_CURRENT_STATE),
		Message: cms,
	})
	return nil
}

func (controller *controller) initModules(ctx context.Context) error {
	for id, mod := range controller.modules {
		if ctx.Err() != nil {
			return nil
		}
		mod.deps = &modutil.ModuleDeps{
			Bus:       controller.bus,
			KV:        *controller.kv.WithPrefix(mod.kvPrefix),
			Log:       controller.log.With("module", id),
			CachePath: filepath.Join(controller.cachePath, id),
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

		b, err := proto.Marshal(&ChangeModuleState{
			ModuleId:    id,
			ModuleState: ModuleState_STARTED,
		})
		if err != nil {
			return fmt.Errorf("marshalling start message for %s: %w", id, err)
		}
		controller.bus.Send(&bus.BusMessage{
			Topic:   BusTopics_CONTROL.String(),
			Type:    int32(MessageType_TYPE_CHANGE_STATE),
			Message: b,
		})
	}
	return nil
}
