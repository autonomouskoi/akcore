package web

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/autonomouskoi/akcore"
	"github.com/autonomouskoi/akcore/bus"
)

type internalHandler struct {
	lock         sync.Mutex
	sendToClient func(*bus.BusMessage) error
	bus          *bus.Bus
	subs         map[string]chan *bus.BusMessage
}

func (h *internalHandler) Close() {
	h.lock.Lock()
	defer h.lock.Unlock()
	for topic, in := range h.subs {
		h.bus.Unsubscribe(topic, in)
	}
	h.bus = nil
}

func (h *internalHandler) handleInternal(msg *bus.BusMessage) error {
	var resp *bus.BusMessage
	var err error
	switch msg.Type {
	case int32(bus.ExternalMessageType_HAS_TOPIC):
		resp, err = handleWith(msg, h.hasTopic)
	case int32(bus.ExternalMessageType_SUBSCRIBE):
		resp, err = handleWith(msg, h.subscribe)
	case int32(bus.ExternalMessageType_UNSUBSCRIBE):
		resp, err = handleWith(msg, h.unsubscribe)
	default:
		return fmt.Errorf("unhandled type %v", msg.Type)
	}
	if err != nil {
		return err
	}
	if resp != nil {
		return h.sendToClient(resp)
	}
	return nil
}

func handleWith[M any, PM akcore.ProtoMessagePointer[M]](
	msg *bus.BusMessage,
	fn func(PM) (protoreflect.ProtoMessage, error),
) (*bus.BusMessage, error) {
	var payload M
	if err := proto.Unmarshal(msg.GetMessage(), PM(&payload)); err != nil {
		return nil, fmt.Errorf("unmarshalling %T proto: %w", payload, err)
	}
	resp, err := fn(&payload)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	respB, err := proto.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshalling response: %w", err)
	}
	return &bus.BusMessage{
		Type:    msg.Type,
		Message: respB,
		ReplyTo: msg.ReplyTo,
	}, nil
}

func (h *internalHandler) hasTopic(payload *bus.HasTopicRequest) (protoreflect.ProtoMessage, error) {
	if payload.TimeoutMs == 0 {
		return &bus.HasTopicResponse{
			Topic:    payload.GetTopic(),
			HasTopic: h.bus.HasTopic(payload.GetTopic()),
		}, nil
	}
	timeout := time.Millisecond * time.Duration(payload.TimeoutMs)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp := &bus.HasTopicResponse{
		Topic: payload.Topic,
	}
	if err := h.bus.WaitForTopic(ctx, payload.Topic, time.Millisecond*10); err == nil {
		resp.HasTopic = true
	}
	return resp, nil
}

func (h *internalHandler) subscribe(payload *bus.SubscribeRequest) (protoreflect.ProtoMessage, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.bus == nil {
		return nil, nil
	}
	if payload.Topic == "" {
		return nil, nil
	}
	in := make(chan *bus.BusMessage, 16)
	h.subs[payload.Topic] = in
	h.bus.Subscribe(payload.Topic, in)
	go func() {
		for msg := range in {
			h.sendToClient(msg)
		}
	}()
	return nil, nil
}

func (h *internalHandler) unsubscribe(payload *bus.UnsubscribeRequest) (protoreflect.ProtoMessage, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	in, present := h.subs[payload.Topic]
	if !present {
		return nil, nil
	}
	h.bus.Unsubscribe(payload.Topic, in)
	return nil, nil
}
