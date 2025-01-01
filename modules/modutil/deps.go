package modutil

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/storage/kv"
)

type ModuleDeps struct {
	Bus         *bus.Bus
	KV          kv.KVPrefix
	Log         *slog.Logger
	StoragePath string
	CachePath   string
}

type Module interface {
	Start(context.Context, *ModuleDeps) error
}

type Deps struct {
	Bus         *bus.Bus
	KV          kv.KV
	Log         *slog.Logger
	Web         Web
	StoragePath string
	CachePath   string
}

type Web interface {
	Handle(path string, handler http.Handler)
}

type ModuleBase struct {
	Log *slog.Logger
	eg  errgroup.Group
}

func (mb *ModuleBase) MarshalMessage(msg *bus.BusMessage, v proto.Message) {
	var err error
	msg.Message, err = proto.Marshal(v)
	if err != nil {
		name := v.ProtoReflect().Descriptor().FullName()
		mb.Log.Error("marshalling", "type", name, "error", err.Error())
		msg.Error = &bus.Error{
			Code:   int32(bus.CommonErrorCode_INVALID_TYPE),
			Detail: proto.String("marshalling " + err.Error()),
		}
	}
}

func (mb *ModuleBase) UnmarshalMessage(msg *bus.BusMessage, v protoreflect.ProtoMessage) *bus.Error {
	if err := proto.Unmarshal(msg.GetMessage(), v); err != nil {
		name := v.ProtoReflect().Descriptor().FullName()
		mb.Log.Error("unmarshalling", "type", name, "error", err.Error())
		return &bus.Error{
			Code:   int32(bus.CommonErrorCode_INVALID_TYPE),
			Detail: proto.String("unmarshalling: " + err.Error()),
		}
	}
	return nil
}

func (mb *ModuleBase) Go(fn func() error) {
	mb.eg.Go(func() error {
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					b := make([]byte, 4096)
					n := runtime.Stack(b, false)
					mb.Log.Error("panic", "panic", r, "stack", string(b[:n]))
					err = fmt.Errorf("panic: %v", r)
				}
			}()
			err = fn()
		}()
		return err
	})
}

func (mb *ModuleBase) Wait() error {
	return mb.eg.Wait()
}
