package modules

import (
	"context"

	"github.com/autonomouskoi/akcore/bus"
)

func (controller *controller) handleRequests(ctx context.Context) error {
	controller.bus.HandleTypes(ctx, BusTopics_MODULE_REQUEST.String(), 8,
		map[int32]bus.MessageHandler{
			int32(MessageTypeRequest_MODULES_LIST_REQ): controller.handleModuleListRequest,
		},
		nil,
	)
	return nil
}

func (controller *controller) handleModuleListRequest(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	controller.lock.Lock()
	defer controller.lock.Unlock()

	mlr := &ModulesListResponse{}
	for _, module := range controller.modules {
		mlr.Entries = append(mlr.Entries, &ModuleListEntry{
			Manifest: module.manifest,
			State: &CurrentModuleState{
				ModuleId:    module.manifest.Id,
				ModuleState: module.state,
				Config:      module.config,
				StateDetail: module.stateDetail,
			},
		})
	}
	controller.MarshalMessage(reply, mlr)
	return reply
}
