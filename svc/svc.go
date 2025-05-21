package svc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	pb "github.com/autonomouskoi/akcore/svc/pb"
	"github.com/autonomouskoi/akcore/svc/webclient"
)

type request struct {
	msg      *bus.BusMessage
	replyVia chan<- *bus.BusMessage
}

type Service struct {
	modutil.ModuleBase
	bus  *bus.Bus
	in   chan request
	lock sync.Mutex
	wc   *webclient.WebClient
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
	deps.Web.Handle(webclientPath, svc.wc)

	return svc, nil
}

func (svc *Service) Start(ctx context.Context) error {
	svc.in = make(chan request, 64)
	svc.Go(svc.handle)

	<-ctx.Done()
	svc.lock.Lock()
	close(svc.in)
	svc.in = nil
	svc.lock.Unlock()

	return svc.Wait()
}

func (svc *Service) Handle(msg *bus.BusMessage) *bus.BusMessage {
	svc.lock.Lock()
	if svc.in == nil {
		svc.lock.Unlock()
		return nil
	}

	replyVia := make(chan *bus.BusMessage)
	svc.in <- request{
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
			default:
				svc.Log.Error("unhandled message type",
					"type", request.msg.GetType(),
					"from_mod", request.msg.GetFromMod(),
				)
				reply = &bus.BusMessage{
					Error: &bus.Error{
						Code: int32(bus.CommonErrorCode_INVALID_TYPE),
					},
				}
			}
			request.replyVia <- reply
		}()
	}
	return nil
}
