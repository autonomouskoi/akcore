package time

import (
	"context"
	"sync"
	"time"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	svc "github.com/autonomouskoi/akcore/svc/pb"
	"golang.org/x/exp/maps"
)

type event struct {
	caller      string
	token       int64
	notifyAfter int64
	repeatEvery int64
}

type Time struct {
	modutil.ModuleBase
	bus       *bus.Bus
	lock      sync.Mutex
	nextToken int64
	events    map[int64]*event
}

func New(deps *modutil.Deps) *Time {
	t := &Time{
		events: map[int64]*event{},
		bus:    deps.Bus,
	}
	t.Log = deps.Log.NewForSource("svc.time")
	return t
}

func (t *Time) HandleNotifyRequest(msg *bus.BusMessage) *bus.BusMessage {
	reply := bus.DefaultReply(msg)
	req := &svc.TimeNotifyRequest{}
	if reply.Error = t.UnmarshalMessage(msg, req); reply.Error != nil {
		return reply
	}

	t.lock.Lock()
	defer t.lock.Unlock()
	token := t.nextToken
	t.nextToken++
	event := &event{
		caller: msg.FromMod,
		token:  token,
	}
	now := time.Now().UnixMilli()
	switch req.TimerType.(type) {
	case *svc.TimeNotifyRequest_At:
		event.notifyAfter = int64(req.GetAt())
	case *svc.TimeNotifyRequest_After:
		event.notifyAfter = now + int64(req.GetAfter())
	case *svc.TimeNotifyRequest_Every:
		event.notifyAfter = now + int64(req.GetEvery())
		event.repeatEvery = int64(req.GetEvery())
	default:
		reply.Error = &bus.Error{
			Code: int32(bus.CommonErrorCode_INVALID_TYPE),
		}
		return reply
	}
	t.events[token] = event
	t.MarshalMessage(reply, &svc.TimeNotifyResponse{Token: token})
	return reply
}

func (t *Time) HandleStopNotificationRequest(msg *bus.BusMessage) *bus.BusMessage {
	reply := bus.DefaultReply(msg)
	req := &svc.TimeStopNotifyRequest{}
	if reply.Error = t.UnmarshalMessage(msg, req); reply.Error != nil {
		return reply
	}

	t.lock.Lock()
	defer t.lock.Unlock()
	for id, event := range t.events {
		if event.caller == msg.FromMod && id == req.GetToken() {
			delete(t.events, id)
			break
		}
	}
	return reply
}

func (t *Time) HandleCurrentTimeRequest(msg *bus.BusMessage) *bus.BusMessage {
	reply := bus.DefaultReply(msg)
	now := time.Now()
	_, offset := now.Zone()
	t.MarshalMessage(reply, &svc.CurrentTimeResponse{
		CurrentTimeMillis: now.UnixMilli(),
		TzOffsetSeconds:   int64(offset),
	})
	return reply
}

func (t *Time) NotifyLoop(ctx context.Context) {
	defer maps.Clear(t.events)
	timer := time.NewTicker(time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-timer.C:
			t.notifyReady(now.UnixMilli())
		}
	}
}

func (t *Time) notifyReady(timeMillis int64) {
	t.lock.Lock()
	defer t.lock.Unlock()
	// send pending events first
	for id, event := range t.events {
		if event.notifyAfter <= timeMillis {
			t.notify(event, timeMillis)
			if event.repeatEvery > 0 {
				t.events[id].notifyAfter = timeMillis + event.repeatEvery // cron-like
			} else {
				delete(t.events, id) // one-time, delivered
			}
		}
	}
	// cleanup events for dead callers
	for id, event := range t.events {
		if !t.bus.HasTopic(event.caller) {
			delete(t.events, id)
			t.Log.Debug("cleaned up event for missing caller",
				"caller", event.caller,
				"token", id,
			)
		}
	}
}

func (t *Time) notify(event *event, timeMillis int64) {
	msg := &bus.BusMessage{
		Topic: event.caller,
		Type:  int32(svc.MessageType_TIME_NOTIFICATION_EVENT),
	}
	if t.MarshalMessage(msg, &svc.TimeNotification{
		Token:             event.token,
		CurrentTimeMillis: timeMillis,
	}); msg.Error != nil {
		return
	}
	t.bus.Send(msg)
}
