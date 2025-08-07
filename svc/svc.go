package svc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	pb "github.com/autonomouskoi/akcore/svc/pb"
	"github.com/autonomouskoi/akcore/svc/template"
	timesvc "github.com/autonomouskoi/akcore/svc/time"
	"github.com/autonomouskoi/akcore/svc/webclient"
)

type request struct {
	ctx      *modutil.PluginContext
	msg      *bus.BusMessage
	replyVia chan<- *bus.BusMessage
}

type Service struct {
	modutil.ModuleBase
	bus  *bus.Bus
	in   chan request
	lock sync.Mutex
	wc   *webclient.WebClient
	time *timesvc.Time
	tmpl *template.Template
}

func New(deps *modutil.Deps) (*Service, error) {
	svc := &Service{}
	svc.Log = deps.Log.NewForSource("svc")
	svc.bus = deps.Bus

	webclientPath := "/s/webclient/"
	wc, err := webclient.New(deps, webclientPath)
	if err != nil {
		return nil, fmt.Errorf("creating webclient: %w", err)
	}
	svc.wc = wc

	tmpl, err := template.New(deps)
	if err != nil {
		return nil, fmt.Errorf("creating template: %w", err)
	}
	svc.tmpl = tmpl

	deps.Web.Handle(webclientPath, svc.wc)

	svc.time = timesvc.New(deps)

	return svc, nil
}

func (svc *Service) Start(ctx context.Context) error {
	svc.in = make(chan request, 64)
	svc.Go(svc.handle)
	svc.Go(func() error { svc.time.NotifyLoop(ctx); return nil })

	<-ctx.Done()
	svc.lock.Lock()
	close(svc.in)
	svc.in = nil
	svc.lock.Unlock()

	return svc.Wait()
}

func (svc *Service) CloseModule(moduleID string) {
	svc.time.CloseModule(moduleID)
	svc.wc.CloseModule(moduleID)
}

func (svc *Service) Handle(pCtx *modutil.PluginContext, msg *bus.BusMessage) *bus.BusMessage {
	svc.lock.Lock()
	if svc.in == nil {
		svc.lock.Unlock()
		return nil
	}

	replyVia := make(chan *bus.BusMessage)
	svc.in <- request{
		ctx:      pCtx,
		msg:      msg,
		replyVia: replyVia,
	}
	svc.lock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var reply *bus.BusMessage
	select {
	case <-ctx.Done():
		reply = &bus.BusMessage{
			Error: &bus.Error{
				Code: int32(bus.CommonErrorCode_TIMEOUT),
			},
		}
	case reply = <-replyVia:
	}
	return reply
}

func (svc *Service) handle() error {
	in := svc.in
	for request := range in {
		func() {
			defer close(request.replyVia)
			var reply *bus.BusMessage
			switch request.msg.GetType() {
			case int32(pb.MessageType_WEBCLIENT_STATIC_DOWNLOAD_REQ):
				reply = svc.wc.HandleRequestStaticDownload(request.msg)
			case int32(pb.MessageType_TEMPLATE_RENDER_REQ):
				reply = svc.tmpl.HandleRequestRender(request.msg)
			case int32(pb.MessageType_TIME_NOTIFICATION_REQ):
				reply = svc.time.HandleNotifyRequest(request.msg)
			case int32(pb.MessageType_TIME_STOP_NOTIFICATION_REQ):
				reply = svc.time.HandleStopNotificationRequest(request.msg)
			case int32(pb.MessageType_TIME_CURRENT_REQ):
				reply = svc.time.HandleCurrentTimeRequest(request.msg)
			case int32(pb.MessageType_WEBCLIENT_HTTP_REQ):
				reply = svc.wc.HandleRequest(*request.ctx, request.msg)
			default:
				svc.Log.Error("unhandled message type",
					"type", request.msg.GetType(),
					"from_mod", request.msg.GetFromMod(),
				)
				reply = &bus.BusMessage{
					Error: &bus.Error{
						Code:   int32(bus.CommonErrorCode_INVALID_TYPE),
						Detail: proto.String("unhandled service call type"),
					},
				}
			}
			request.replyVia <- reply
		}()
	}
	return nil
}
