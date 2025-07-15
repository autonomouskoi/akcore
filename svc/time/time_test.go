package time_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/akcore/svc/log"
	svc "github.com/autonomouskoi/akcore/svc/pb"
	timesvc "github.com/autonomouskoi/akcore/svc/time"
)

func TestTime(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logDir := t.TempDir()
	ml, err := log.New(logDir, &svc.Config{})
	require.NoError(t, err, "creating master logger")
	defer ml.Close()

	deps := &modutil.Deps{
		Bus: bus.New(ctx),
		Log: ml,
	}

	ts := timesvc.New(deps)
	go ts.NotifyLoop(ctx)

	modID := "mod-id"
	in := make(chan *bus.BusMessage)
	deps.Bus.Subscribe(modID, in)

	t.Run("at", func(t *testing.T) {
		msg := &bus.BusMessage{
			Type:    int32(svc.MessageType_TIME_NOTIFICATION_REQ),
			FromMod: modID,
		}
		req := &svc.TimeNotifyRequest{
			TimerType: &svc.TimeNotifyRequest_At{
				At: uint64(time.Now().UnixMilli() + 5),
			},
		}
		var err error
		msg.Message, err = proto.Marshal(req)
		require.NoError(t, err, "marshalling request")

		reply := ts.HandleNotifyRequest(msg)
		require.Nil(t, reply.Error, "in reply")

		resp := &svc.TimeNotifyResponse{}
		require.NoError(t, proto.Unmarshal(reply.GetMessage(), resp), "unmarshalling response")
		require.Zero(t, resp.Token)

		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		select {
		case msg = <-in:
			// cool
		case <-ctx.Done():
			t.Fatal("waiting for notification")
		}
		require.Equal(t, int32(svc.MessageType_TIME_NOTIFICATION_EVENT), msg.GetType())
		notif := &svc.TimeNotification{}
		require.NoError(t, proto.Unmarshal(msg.GetMessage(), notif), "unmarshalling notification")

		now := time.Now().UnixMilli()
		require.LessOrEqual(t, notif.GetCurrentTimeMillis(), now+5)
		require.GreaterOrEqual(t, notif.GetCurrentTimeMillis(), now-5)
		require.Equal(t, resp.GetToken(), notif.GetToken())
	})

	t.Run("after", func(t *testing.T) {
		msg := &bus.BusMessage{
			Type:    int32(svc.MessageType_TIME_NOTIFICATION_REQ),
			FromMod: modID,
		}
		req := &svc.TimeNotifyRequest{
			TimerType: &svc.TimeNotifyRequest_After{
				After: 5,
			},
		}
		var err error
		msg.Message, err = proto.Marshal(req)
		require.NoError(t, err, "marshalling request")

		reply := ts.HandleNotifyRequest(msg)
		require.Nil(t, reply.Error, "in reply")

		resp := &svc.TimeNotifyResponse{}
		require.NoError(t, proto.Unmarshal(reply.GetMessage(), resp), "unmarshalling response")
		require.NotZero(t, resp.GetToken())

		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		select {
		case msg = <-in:
			// cool
		case <-ctx.Done():
			t.Fatal("waiting for notification")
		}
		require.Equal(t, int32(svc.MessageType_TIME_NOTIFICATION_EVENT), msg.GetType())
		notif := &svc.TimeNotification{}
		require.NoError(t, proto.Unmarshal(msg.GetMessage(), notif), "unmarshalling notification")

		now := time.Now().UnixMilli()
		require.LessOrEqual(t, notif.GetCurrentTimeMillis(), now+5)
		require.GreaterOrEqual(t, notif.GetCurrentTimeMillis(), now-5)
		require.Equal(t, resp.GetToken(), notif.GetToken())
	})

	t.Run("every", func(t *testing.T) {
		msg := &bus.BusMessage{
			Type:    int32(svc.MessageType_TIME_NOTIFICATION_REQ),
			FromMod: modID,
		}
		req := &svc.TimeNotifyRequest{
			TimerType: &svc.TimeNotifyRequest_Every{
				Every: 5,
			},
		}
		var err error
		msg.Message, err = proto.Marshal(req)
		require.NoError(t, err, "marshalling request")

		reply := ts.HandleNotifyRequest(msg)
		require.Nil(t, reply.Error, "in reply")

		resp := &svc.TimeNotifyResponse{}
		require.NoError(t, proto.Unmarshal(reply.GetMessage(), resp), "unmarshalling response")
		require.NotZero(t, resp.GetToken())

		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		for range 3 {
			select {
			case msg = <-in:
				// cool
			case <-ctx.Done():
				t.Fatal("waiting for notification")
			}
			require.Equal(t, int32(svc.MessageType_TIME_NOTIFICATION_EVENT), msg.GetType())
			notif := &svc.TimeNotification{}
			require.NoError(t, proto.Unmarshal(msg.GetMessage(), notif), "unmarshalling notification")

			now := time.Now().UnixMilli()
			require.LessOrEqual(t, notif.GetCurrentTimeMillis(), now+5)
			require.GreaterOrEqual(t, notif.GetCurrentTimeMillis(), now-5)
			require.Equal(t, resp.GetToken(), notif.GetToken())
		}

		// cancel
		msg = &bus.BusMessage{
			Type:    int32(svc.MessageType_TIME_STOP_NOTIFICATION_REQ),
			FromMod: modID,
		}
		msg.Message, err = proto.Marshal(&svc.TimeStopNotifyRequest{Token: resp.GetToken()})
		require.NoError(t, err, "marshalling stop reequest")
		reply = ts.HandleStopNotificationRequest(msg)
		require.Nil(t, reply.Error, "in stop reply")

		// there should be at most 1 more notification
		received := 0
		cleanupCtx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		defer cancel()
	CLEANUPLOOP:
		for {
			select {
			case <-in:
				received++
			case <-cleanupCtx.Done():
				break CLEANUPLOOP
			}
		}
		require.Less(t, received, 2)
	})
}
