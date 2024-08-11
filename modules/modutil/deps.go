package modutil

import (
	"context"
	"log/slog"
	"net/http"

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
