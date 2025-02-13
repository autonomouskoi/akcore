package bustest

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/akcore/storage/kv"
)

func NewDeps(t *testing.T) (context.Context, func(), *modutil.ModuleDeps) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	t.Cleanup(cancel)

	b := bus.New(ctx)
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	kvm, err := kv.NewMemory()
	require.NoError(t, err, "creating memory kv")

	deps := &modutil.ModuleDeps{
		Bus: b,
		Log: log,
		KV:  *kvm.WithPrefix([8]byte{0, 1, 2, 3, 4, 5, 6, 7}),
	}

	return ctx, cancel, deps
}

func AssertPayload(t *testing.T, msg *bus.BusMessage, v proto.Message) {
	t.Helper()

	if err := proto.Unmarshal(msg.GetMessage(), v); err != nil {
		t.Fatalf("unmarshalling message of type %T: %v", v, err)
	}
}
