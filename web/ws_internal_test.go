package web

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"nhooyr.io/websocket"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/bus/bustest"
)

func marshaller(t *testing.T) func(p protoreflect.ProtoMessage) []byte {
	return func(p protoreflect.ProtoMessage) []byte {
		t.Helper()
		b, err := proto.Marshal(p)
		require.NoError(t, err, "marshalling")
		return b
	}
}

func unmarshaller(t *testing.T) func(b []byte, p protoreflect.ProtoMessage) {
	return func(b []byte, p protoreflect.ProtoMessage) {
		t.Helper()
		require.NoError(t, proto.Unmarshal(b, p), "unmarshalling")
	}
}

func sender(t *testing.T, wc *websocket.Conn) func(context.Context, *bus.BusMessage) {
	return func(ctx context.Context, bm *bus.BusMessage) {
		t.Helper()
		b, err := proto.Marshal(bm)
		require.NoError(t, err, "marshalling BusMessage")
		require.NoError(t, wc.Write(ctx, websocket.MessageBinary, b), "sending")
	}
}

func reader(t *testing.T, wc *websocket.Conn) func(context.Context) *bus.BusMessage {
	return func(ctx context.Context) *bus.BusMessage {
		t.Helper()
		typ, b, err := wc.Read(ctx)
		require.NoError(t, err, "reading BusMessage")
		require.Equal(t, websocket.MessageBinary, typ)
		bm := &bus.BusMessage{}
		require.NoError(t, proto.Unmarshal(b, bm), "unmarshalling BusMessage")
		return bm
	}
}

func TestHasTopic(t *testing.T) {
	t.Parallel()
	requireT := require.New(t)

	ctx, _, deps := bustest.NewDeps(t)

	h := newWS(deps)
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	c := server.Client()

	wc, _, err := websocket.Dial(ctx, server.URL, &websocket.DialOptions{
		HTTPClient: c,
	})
	requireT.NoError(err, "dialing")
	t.Cleanup(func() { wc.CloseNow() })

	marshal := marshaller(t)
	unmarshal := unmarshaller(t)
	sendBM := sender(t, wc)
	readBM := reader(t, wc)

	sendBM(ctx, &bus.BusMessage{
		Type: int32(bus.ExternalMessageType_HAS_TOPIC),
		Message: marshal(&bus.HasTopicRequest{
			Topic: "foo",
		}),
		ReplyTo: proto.Int64(123),
	})

	bm := readBM(ctx)
	requireT.Equal(int32(bus.ExternalMessageType_HAS_TOPIC), bm.Type)
	requireT.Equal(int64(123), bm.GetReplyTo())

	resp := &bus.HasTopicResponse{}
	unmarshal(bm.GetMessage(), resp)
	requireT.False(resp.HasTopic)

	in := make(chan *bus.BusMessage)
	deps.Bus.Subscribe("foo", in)

	sendBM(ctx, &bus.BusMessage{
		Type: int32(bus.ExternalMessageType_HAS_TOPIC),
		Message: marshal(&bus.HasTopicRequest{
			Topic: "foo",
		}),
		ReplyTo: proto.Int64(124),
	})

	bm = readBM(ctx)
	requireT.Equal(int32(bus.ExternalMessageType_HAS_TOPIC), bm.Type)
	requireT.Equal(int64(124), bm.GetReplyTo())

	resp = &bus.HasTopicResponse{}
	unmarshal(bm.GetMessage(), resp)
	requireT.True(resp.HasTopic)
}

func TestSubscribe(t *testing.T) {
	t.Parallel()
	requireT := require.New(t)

	ctx, _, deps := bustest.NewDeps(t)

	h := newWS(deps)
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	c := server.Client()

	wc, _, err := websocket.Dial(ctx, server.URL, &websocket.DialOptions{
		HTTPClient: c,
	})
	requireT.NoError(err, "dialing")
	t.Cleanup(func() { wc.CloseNow() })

	readBM := reader(t, wc)
	marshal := marshaller(t)
	sendBM := sender(t, wc)

	sendBM(ctx, &bus.BusMessage{
		Type: int32(bus.ExternalMessageType_SUBSCRIBE),
		Message: marshal(&bus.SubscribeRequest{
			Topic: "test-topic",
		}),
	})

	// wait a moment for the subscribe to process
	time.Sleep(time.Millisecond * 10)
	deps.Bus.Send(&bus.BusMessage{
		Topic: "test-topic",
	})

	bm := readBM(ctx)
	requireT.Equal("test-topic", bm.Topic)

	sendBM(ctx, &bus.BusMessage{
		Type: int32(bus.ExternalMessageType_UNSUBSCRIBE),
		Message: marshal(&bus.UnsubscribeRequest{
			Topic: "test-topic",
		}),
	})

	// wait a moment for the unsubscribe to process
	time.Sleep(time.Millisecond * 10)
	deps.Bus.Send(&bus.BusMessage{
		Topic: "test-topic",
	})

	readCtx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	t.Cleanup(cancel)
	_, _, err = wc.Read(readCtx)
	requireT.ErrorIs(err, context.DeadlineExceeded)
}
