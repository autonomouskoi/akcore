package modules

import (
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"

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
		ModuleId: m.manifest.Id,
		//ModuleName:  m.manifest.Name,
		ModuleState: m.state,
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
