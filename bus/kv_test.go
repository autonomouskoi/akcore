package bus_test

import (
	"context"
	"testing"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/bus/bustest"
	"github.com/autonomouskoi/akcore/modules"
	"github.com/autonomouskoi/akcore/modules/modutil"

	"github.com/stretchr/testify/require"
)

func NewKVDeps(t *testing.T) (context.Context, func(), *modules.WASM, *modutil.ModuleDeps) {
	t.Helper()
	ctx, cancel, deps := bustest.NewDeps(t)

	w := &modules.WASM{}
	w.SetWASMFiles("tinygo/wasm_out/test.wasm")
	return ctx, cancel, w, deps
}

func TestKV(t *testing.T) {
	t.Parallel()

	ctx, cancel, w, deps := NewKVDeps(t)

	go func() {
		require.NoError(t, w.Start(ctx, deps))
		cancel()
	}()

	in := make(chan *bus.BusMessage)
	deps.Bus.Subscribe("HOST", in)
	for msg := range in {
		if msg.Error == nil {
			t.Logf("%s PASS", string(msg.GetMessage()))
		} else {
			t.Errorf("%s FAIL: %s", string(msg.GetMessage()), msg.GetError().GetDetail())
		}
	}
}
