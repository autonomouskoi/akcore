package bus_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/bus/bustest"
)

func TestSubAndSend(t *testing.T) {
	t.Parallel()
	requireT := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	b := bus.New(ctx)

	clientA := make(chan *bus.BusMessage, 5)
	clientB := make(chan *bus.BusMessage, 5)
	clientC := make(chan *bus.BusMessage, 5)

	b.Subscribe("test", clientA)
	b.Subscribe("test", clientB)
	b.Subscribe("other", clientC)

	m := &bus.BusMessage{
		Topic: "test",
	}

	b.Send(m)

	select {
	case received := <-clientA:
		requireT.Equal(m, received)
	default:
		t.Fatal("no value for clientA")
	}
	select {
	case received := <-clientB:
		requireT.Equal(m, received)
	default:
		t.Fatal("no value for clientB")
	}
	select {
	case <-clientC:
		t.Fatal("unexpected value for clientC")
	default:
	}

	m2 := &bus.BusMessage{
		Topic: "test",
	}
	// test discarding. Send more than the recipient channels can hold
	for i := 0; i < 10; i++ {
		b.Send(m2)
	}

	aCount := 0
	bCount := 0
SELECTLOOP:
	for {
		select {
		case received := <-clientA:
			requireT.Equal(m2, received)
			aCount++
		case received := <-clientB:
			requireT.Equal(m2, received)
			bCount++
		default:
			break SELECTLOOP
		}
	}
	requireT.Equal(5, aCount)
	requireT.Equal(5, bCount)
}

func TestUnsubscribe(t *testing.T) {
	t.Parallel()
	requireT := require.New(t)

	_, _, deps := bustest.NewDeps(t)
	b := deps.Bus

	topic := "test-topic"
	in := make(chan *bus.BusMessage, 4)
	b.Subscribe(topic, in)

	requireT.True(b.HasTopic(topic))
	b.Send(&bus.BusMessage{Topic: topic, Type: 17})
	msg := <-in
	requireT.NotNil(msg)
	requireT.Equal(int32(17), msg.Type)

	requireT.True(b.Unsubscribe(topic, in))
	requireT.False(b.HasTopic(topic))
	requireT.False(b.Unsubscribe(topic, in))

	// the channel isn't closed
	in <- &bus.BusMessage{Type: 23}
	msg = <-in
	requireT.NotNil(msg)
	requireT.Equal(int32(23), msg.Type)

	close(in)
}
