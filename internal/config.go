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
)

var (
	cfgKey = []byte("config")
)

type service struct {
	modutil.ModuleBase
	lock sync.Mutex
	bus  *bus.Bus
	kv   kv.KVPrefix
	cfg  *Config
}

func Start(ctx context.Context, deps *modutil.Deps) error {
	svc := &service{
		bus: deps.Bus,
		kv:  *deps.KV.WithPrefix([8]byte{}),
	}
	svc.Log = deps.Log.With("module", "internal")
	var err error
	svc.cfg, err = getConfig(svc.kv)
	if err != nil {
		return fmt.Errorf("getting config: %w", err)
	}
	eg := &errgroup.Group{}

	eg.Go(func() error { return svc.handleRequests(ctx) })
	eg.Go(func() error { return svc.handleCommands(ctx) })

	return eg.Wait()
}

func (svc *service) handleRequests(ctx context.Context) error {
	svc.bus.HandleTypes(ctx, BusTopic_INTERNAL_REQUEST.String(), 8,
		map[int32]bus.MessageHandler{
			int32(MessageTypeRequest_CONFIG_GET_REQ): svc.handleRequestConfigGet,
		},
		nil,
	)
	return nil
}

func (svc *service) handleRequestConfigGet(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	svc.lock.Lock()
	svc.MarshalMessage(reply, &ConfigGetResponse{Config: svc.cfg})
	svc.lock.Unlock()
	return reply
}

func (svc *service) handleCommands(ctx context.Context) error {
	svc.bus.HandleTypes(ctx, BusTopic_INTERNAL_COMMAND.String(), 4,
		map[int32]bus.MessageHandler{
			int32(MessageTypeCommand_CONFIG_SET_REQ): svc.handleCommandConfigSet,
		},
		nil,
	)
	return nil
}

func (svc *service) handleCommandConfigSet(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	csr := &ConfigSetRequest{}
	if reply.Error = svc.UnmarshalMessage(msg, csr); reply.Error != nil {
		return reply
	}
	svc.lock.Lock()
	svc.cfg = csr.GetConfig()
	if err := setConfig(svc.kv, svc.cfg); err != nil {
		reply.Error = &bus.Error{
			Detail: proto.String(err.Error()),
		}
		svc.Log.Error("saving config", "error", err.Error())
		return reply
	}
	svc.MarshalMessage(reply, &ConfigSetResponse{Config: svc.cfg})
	svc.lock.Unlock()
	return reply
}

func getConfig(kv kv.KVPrefix) (*Config, error) {
	cfg := &Config{}
	if err := kv.GetProto(cfgKey, cfg); err != nil && !errors.Is(err, akcore.ErrNotFound) {
		return nil, err
	}
	return cfg, nil
}

func setConfig(kv kv.KVPrefix, cfg *Config) error {
	return kv.SetProto(cfgKey, cfg)
}
