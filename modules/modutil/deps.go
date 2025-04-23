package modutil

import (
	"context"
	_ "embed"
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

type Service interface {
	Handle(msg *bus.BusMessage) *bus.BusMessage
}

// ModuleDeps carries the deps specific to a module.
type ModuleDeps struct {
	Bus         *bus.Bus
	KV          kv.KVPrefix
	Log         *slog.Logger
	StoragePath string
	CachePath   string
	Svc         Service
}

// A Module can be started with context and deps
type Module interface {
	Start(context.Context, *ModuleDeps) error
	Icon() ([]byte, string, error)
}

// Deps carries the deps of the modules system itself, not deps for a specific
// module
type Deps struct {
	Bus         *bus.Bus
	KV          kv.KV
	Log         *slog.Logger
	Web         Web
	StoragePath string
	CachePath   string
	HttpClient  *http.Client
}

// Web things can handle HTTP requests
type Web interface {
	Handle(path string, handler http.Handler)
}

// ModuleBase provides some common module functionality
type ModuleBase struct {
	Log *slog.Logger
	eg  errgroup.Group
}

// MarshalMessage marshals a payload and sets it on the provided BusMessage. If
// marshalling fails, an error is logged and msg.Error is set
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

// UnmarshalMessage unmarshals a payload from a BusMessage. If unmarshalling
// fails, an error is logged and a *bus.Error is returned. A useful idiom is:
//
//	if reply.Error = m.UnmarshallMessage(msg, target); reply.Error != nil {
//	    return reply
//	}
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

// Go launches a goroutine with the provided function using the internal
// errgroup
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

// Wait for the internal errgroup to finish.
func (mb *ModuleBase) Wait() error {
	return mb.eg.Wait()
}

//go:embed ak_logo.svg
var AKLogo []byte

const AKLogoType = "image/svg+xml"

// Icon returns a default icon and MIME type
func (mb *ModuleBase) Icon() ([]byte, string, error) {
	return AKLogo, AKLogoType, nil
}
