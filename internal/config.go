package internal

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore"
	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/akcore/storage/kv"
	svc "github.com/autonomouskoi/akcore/svc/pb"
)

var (
	cfgKey = []byte("config")
)

// our internal service acts as a module but represents internal functionality
type service struct {
	modutil.ModuleBase
	lock       sync.Mutex
	bus        *bus.Bus
	kv         kv.KVPrefix
	cfg        *svc.Config
	cfgUpdated func(*svc.Config)
}

// Start internal functions
func Start(ctx context.Context, deps *modutil.Deps) error {
	svc := &service{
		bus: deps.Bus,
		kv:  *deps.KV.WithPrefix([8]byte{}),
	}
	svc.Log = deps.Log.NewForSource("internal")
	var err error
	svc.cfg, err = GetConfig(svc.kv)
	if err != nil {
		return fmt.Errorf("getting config: %w", err)
	}
	svc.cfgUpdated = deps.UpdateConfig
	eg := &errgroup.Group{}

	eg.Go(func() error { return svc.handleRequests(ctx) })
	eg.Go(func() error { return svc.handleCommands(ctx) })

	return eg.Wait()
}

// handle messages on the internal request topic
func (s *service) handleRequests(ctx context.Context) error {
	s.bus.HandleTypes(ctx, svc.BusTopic_INTERNAL_REQUEST.String(), 8,
		map[int32]bus.MessageHandler{
			int32(svc.MessageTypeRequest_CONFIG_GET_REQ): s.handleRequestConfigGet,
		},
		nil,
	)
	return nil
}

// handle requests to fetch the config
func (s *service) handleRequestConfigGet(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	s.lock.Lock()
	s.MarshalMessage(reply, &svc.ConfigGetResponse{Config: s.cfg})
	s.lock.Unlock()
	return reply
}

// handle messages on the internal command topic
func (s *service) handleCommands(ctx context.Context) error {
	s.bus.HandleTypes(ctx, svc.BusTopic_INTERNAL_COMMAND.String(), 4,
		map[int32]bus.MessageHandler{
			int32(svc.MessageTypeCommand_CONFIG_SET_REQ): s.handleCommandConfigSet,
		},
		nil,
	)
	return nil
}

// handle requests to update the stored config
func (s *service) handleCommandConfigSet(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	csr := &svc.ConfigSetRequest{}
	if reply.Error = s.UnmarshalMessage(msg, csr); reply.Error != nil {
		return reply
	}
	s.lock.Lock()
	s.cfg = csr.GetConfig()
	if err := setConfig(s.kv, s.cfg); err != nil {
		reply.Error = &bus.Error{
			Detail: proto.String(err.Error()),
		}
		s.Log.Error("saving config", "error", err.Error())
		return reply
	}
	s.MarshalMessage(reply, &svc.ConfigSetResponse{Config: s.cfg})
	s.cfgUpdated(s.cfg)
	s.lock.Unlock()
	return reply
}

func GetConfig(kv kv.KVPrefix) (*svc.Config, error) {
	cfg := &svc.Config{}
	if err := kv.GetProto(cfgKey, cfg); err != nil && !errors.Is(err, akcore.ErrNotFound) {
		return nil, err
	}
	return cfg, nil
}

func setConfig(kv kv.KVPrefix, cfg *svc.Config) error {
	return kv.SetProto(cfgKey, cfg)
}
