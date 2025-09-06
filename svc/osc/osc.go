// Package osc provides Open Sound Control message sending as an internal service
package osc

import (
	"slices"
	"sync"

	"github.com/hypebeast/go-osc/osc"
	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	svc "github.com/autonomouskoi/akcore/svc/pb"
	"github.com/autonomouskoi/datastruct/iterutil"
	"github.com/autonomouskoi/datastruct/mapset"
)

type OSC struct {
	modutil.ModuleBase
	bus     *bus.Bus
	cfg     *svc.OSCConfig
	lock    sync.Mutex
	clients map[string]*osc.Client
}

func New(deps *modutil.Deps) *OSC {
	o := &OSC{
		bus:     deps.Bus,
		clients: map[string]*osc.Client{},
	}
	o.Log = deps.Log.NewForSource("svc.osc")
	o.UpdateConfig(deps.Config)
	deps.WatchConfig(o.UpdateConfig)
	return o
}

func (o *OSC) UpdateConfig(cfg *svc.Config) {
	o.lock.Lock()
	o.cfg = cfg.GetOscConfig()
	wantTargets := mapset.FromSeq(iterutil.Map(
		slices.Values(o.cfg.GetTargets()),
		func(t *svc.OSCTarget) string { return t.GetName() },
	))
	// remove the clients we don't want
	for name := range o.clients {
		if !wantTargets.Has(name) {
			delete(o.clients, name)
		}
	}
	// create the clients we're missing
	for _, target := range o.cfg.GetTargets() {
		if _, present := o.clients[target.Name]; !present {
			o.clients[target.Name] = osc.NewClient(target.GetAddress(), int(target.GetPort()))
		}
	}
	o.lock.Unlock()
}

func (o *OSC) CloseModule(moduleID string) {}

func (o *OSC) HandleRequestListTargets(msg *bus.BusMessage) *bus.BusMessage {
	reply := bus.DefaultReply(msg)

	o.lock.Lock()
	o.MarshalMessage(reply, &svc.OSCListTargetsResponse{
		Targets: o.cfg.GetTargets(),
	})
	o.lock.Unlock()

	return reply
}

func (o *OSC) HandleRequestSendMessage(msg *bus.BusMessage) *bus.BusMessage {
	reply := bus.DefaultReply(msg)
	smr := &svc.OSCSendMessageRequest{}
	if reply.Error = o.UnmarshalMessage(msg, smr); reply.Error != nil {
		return reply
	}
	if reply.Error = o.sendMessage(smr); reply.Error != nil {
		return reply
	}
	o.MarshalMessage(reply, &svc.OSCSendMessageResponse{})

	return reply
}

func (o *OSC) sendMessage(smr *svc.OSCSendMessageRequest) *bus.Error {
	targetName := smr.GetTargetName()
	o.lock.Lock()
	client := o.clients[targetName]
	o.lock.Unlock()
	if client == nil {
		o.Log.Error("invalid target", "target_id", targetName)
		return &bus.Error{
			Code:   int32(bus.CommonErrorCode_INVALID_TYPE),
			Detail: proto.String("invalid target: " + targetName),
		}
	}

	oscMsg := osc.NewMessage(smr.GetAddress())
	for _, oscVal := range smr.GetValues() {
		var value any
		switch v := oscVal.Value.(type) {
		case *svc.OSCValue_Nil:
			// value nil is desired
		case *svc.OSCValue_Int32:
			value = v.Int32
		case *svc.OSCValue_Float32:
			value = v.Float32
		case *svc.OSCValue_String_:
			value = v.String_
		case *svc.OSCValue_Blob:
			value = v.Blob
		case *svc.OSCValue_Int64:
			value = v.Int64
		/*
			case *svc.OSCValue_Time:
				value = v.Time
			case *svc.OSCValue_Double:
				value = v.Double
		*/
		case *svc.OSCValue_True:
			value = true
		case *svc.OSCValue_False:
			value = false
		}
		oscMsg.Append(value)
	}
	if err := client.Send(oscMsg); err != nil {
		o.Log.Error("sending osc",
			"host", client.IP(),
			"port", client.Port(),
			"osc", *oscMsg,
			"error", err.Error(),
		)
		return &bus.Error{
			Code:   int32(bus.CommonErrorCode_UNKNOWN),
			Detail: proto.String("sending: " + err.Error()),
		}
	}
	return nil
}
