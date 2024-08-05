package bus_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonomouskoi/akcore/bus"
)

func TestReply(t *testing.T) {
	t.Parallel()
	requireT := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	b := bus.New(ctx)

	topic := "echo"

	// start the echo service
	echoIn := make(chan *bus.BusMessage, 5)
	b.Subscribe(topic, echoIn)
	go func() {
		for in := range echoIn {
			s := in.Message
			b.SendReply(in, &bus.BusMessage{
				Message: s,
			})
		}
	}()

	ctx, cancel = context.WithTimeout(ctx, time.Second*10)
	t.Cleanup(cancel)

	testValue := []byte("test value")

	reply := b.WaitForReply(ctx,
		&bus.BusMessage{
			Topic:   topic,
			Message: testValue,
		},
	)
	requireT.Nil(reply.Error, "waiting for reply")
	requireT.Equal(testValue, reply.Message)
}
