package modules

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
)

func (controller *controller) handleCommand(ctx context.Context) error {
	controller.bus.HandleTypes(ctx, BusTopics_MODULE_COMMAND.String(), 8,
		map[int32]bus.MessageHandler{
			int32(MessageTypeCommand_MODULE_AUTOSTART_SET_REQ): controller.handleCommandModuleAutostartSet,
			int32(MessageTypeCommand_MODULE_STATE_SET_REQ):     controller.handleCommandModuleStateSet,
		},
		nil,
	)
	return nil
}

func (controller *controller) handleCommandModuleAutostartSet(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	masr := &ModuleAutostartSetRequest{}
	if reply.Error = controller.UnmarshalMessage(msg, masr); reply.Error != nil {
		return reply
	}
	module, ok := controller.modules[masr.ModuleId]
	if !ok {
		reply.Error = &bus.Error{
			Code:   int32(bus.CommonErrorCode_NOT_FOUND),
			Detail: proto.String("not found"),
		}
		return reply
	}
	module.config.AutomaticStart = masr.GetAutostart()
	module.sendState()
	return reply
}

func (controller *controller) handleCommandModuleStateSet(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	mssr := &ModuleStateSetRequest{}
	if reply.Error = controller.UnmarshalMessage(msg, mssr); reply.Error != nil {
		return reply
	}
	switch mssr.State {
	case ModuleState_STARTED:
		controller.startModule(controller.runningCtx, mssr.ModuleId)
	case ModuleState_STOPPED:
		controller.stopModule(mssr.ModuleId)
	default:
		controller.Log.Error("unhandled module state",
			"module_id", mssr.ModuleId, "state", mssr.State)
	}
	return reply
}
