package bustest

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
)

func NewDeps(t *testing.T) (context.Context, func(), *modutil.ModuleDeps) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	t.Cleanup(cancel)

	b := bus.New(ctx)
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	deps := &modutil.ModuleDeps{
		Bus: b,
		Log: log,
	}

	return ctx, cancel, deps
}

func AssertPayload(t *testing.T, msg *bus.BusMessage, v proto.Message) {
	t.Helper()

	if err := proto.Unmarshal(msg.GetMessage(), v); err != nil {
		t.Fatalf("unmarshalling message of type %T: %v", v, err)
	}
}
