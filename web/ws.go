package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"

	"github.com/autonomouskoi/akcore"
	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
)

// WS manages websockets
type WS struct {
	bus *bus.Bus
	log akcore.Logger
}

func newWS(deps *modutil.Deps) *WS {
	ws := &WS{
		bus: deps.Bus,
		log: deps.Log.With("module", "ws"),
	}
	return ws
}

// ServeHTTP handles websocket requests
func (ws *WS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws.log.Debug("handling WS", "url", r.URL.String())

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		ws.log.Error("accepting websocket", "error", err.Error())
		return
	}
	ws.log.Debug("accepted websocket")
	defer c.CloseNow()

	eg, wsCtx := errgroup.WithContext(context.Background())

	// create a channel and function to send messages to the client
	toClient := make(chan *bus.BusMessage, 16)
	sendToClient := func(msg *bus.BusMessage) error {
		msgB, err := proto.Marshal(msg)
		if err != nil {
			ws.log.Error("marshalling message to client",
				"topic", msg.GetTopic(),
				"error", err.Error(),
			)
			return nil
		}
		if err := c.Write(wsCtx, websocket.MessageBinary, msgB); err != nil {
			return fmt.Errorf("writing message: %w", err)
		}
		return nil
	}

	// receive messages bound for the client and send them through the websocket
	eg.Go(func() error {
		defer bus.Drain(toClient)
		for msg := range toClient {
			if err := sendToClient(msg); err != nil {
				return err
			}
		}
		return errors.New("input closed")
	})

	// handle messages from the client
	eg.Go(func() error {
		ih := &internalHandler{
			sendToClient: sendToClient,
			bus:          ws.bus,
			subs:         map[string]chan *bus.BusMessage{},
		}
		defer ih.Close()
		defer close(toClient)
		for {
			// receive a message from the client
			typ, msgB, err := c.Read(wsCtx)
			if err != nil {
				return fmt.Errorf("reading: %w", err)
			}
			if typ != websocket.MessageBinary {
				ws.log.Error("non-binary message received",
					"msg", string(msgB),
				)
				continue
			}
			// unmarshal it to a BusMessage
			msg := new(bus.BusMessage)
			if err := proto.Unmarshal(msgB, msg); err != nil {
				ws.log.Error("unmarshaling BusMessage", "error", err.Error())
				continue
			}
			// if a topic is specified, it's for the bus
			if msg.GetTopic() != "" {
				if msg.GetReplyTo() == 0 {
					ws.bus.Send(msg)
				} else {
					go func() {
						ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
						originalReplyTo := msg.ReplyTo
						defer cancel()
						response := ws.bus.WaitForReply(ctx, msg)
						response.ReplyTo = originalReplyTo
						toClient <- response
					}()
				}
				continue
			}
			go func() {
				// it has no topic, it's a message related to internal
				// functionality
				if err := ih.handleInternal(msg); err != nil {
					ws.log.Error("handling internal message",
						"type", msg.Type,
						"error", err.Error(),
					)
				}
			}()
		}
	})

	if err := eg.Wait(); err != nil {
		ws.log.Error("handling WS", "error", err.Error())
	}
	ws.log.Debug("exiting WS handler")
}
