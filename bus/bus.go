package bus

import (
	"context"
	"sync"
	"time"

	"github.com/autonomouskoi/datastruct/mapset"
)

// Bus delivers messages between modules
type Bus struct {
	lock           sync.Mutex
	subs           map[string]mapset.MapSet[chan<- *BusMessage]
	pendingReplies map[int64]pendingReply
}

// New creates a new Bus. When ctx is cancelled the Bus will close all channels,
// signalling to modules that no more messages will be delivered.
func New(ctx context.Context) *Bus {
	b := &Bus{
		subs:           map[string]mapset.MapSet[chan<- *BusMessage]{},
		pendingReplies: map[int64]pendingReply{},
	}

	go b.prunePendingRepliesEvery(ctx, time.Minute)
	go func() {
		// on cancellation, lock and start closing channels to shut down
		<-ctx.Done()
		b.lock.Lock()
		defer b.lock.Unlock()
		for _, subs := range b.subs {
			for sub := range subs {
				close(sub)
			}
		}
		for _, pr := range b.pendingReplies {
			close(pr.out)
		}
		b.subs = nil // closed for business
		b.pendingReplies = nil
	}()

	return b
}

// HasTopic returns whether or not a topic has any subscribers. Note that
// reality may change immediately after HasTopic returns.
func (b *Bus) HasTopic(topic string) bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return len(b.subs[topic]) > 0
}

// WaitForTopic waits for a topic to have at least one subscriber. This function
// will poll repeatedly, sleeping for checkEvery between polls. If ctx returns
// an error, polling stops and the error returned by ctx.Err is returned.
func (b *Bus) WaitForTopic(ctx context.Context, topic string, checkEvery time.Duration) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if b.HasTopic(topic) {
			return nil
		}
		time.Sleep(checkEvery)
	}
}

// Subscribe to a topic with messages to be received on the provided channel.
// The caller should not close the channel.
func (b *Bus) Subscribe(topic string, recv chan<- *BusMessage) {
	b.lock.Lock()
	defer b.lock.Unlock()
	subs, present := b.subs[topic]
	if !present {
		subs = map[chan<- *BusMessage]struct{}{}
		b.subs[topic] = subs
	}
	subs.Add(recv)
}

// Send a message on the bus. The destination topic will be taken from msg.
func (b *Bus) Send(msg *BusMessage) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for recv := range b.subs[msg.GetTopic()] {
		select {
		case recv <- msg:
			// cool
		default:
			// the recipient can't receive the message. Drop it
		}
	}
}

// Unsubscribe the channel from the topic. Returns whether or not the channel
// was subscribed to the topic. Unsubscribe closes the channel.
func (b *Bus) Unsubscribe(topic string, recv chan<- *BusMessage) bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	if !b.subs[topic].Has(recv) {
		return false
	}
	delete(b.subs[topic], recv)
	close(recv)
	return true
}

// A MessageHandler processes a message and returns a reply or nil
type MessageHandler func(*BusMessage) *BusMessage

// HandleTypes creates a channel with chanCap capacity and subscribes it to
// topic. When ctx is cancelled, the topic is unsubscribed. When a message is
// received on the channel, the handlers map is searched for a handler matching
// the message type. If a match is found that handler is invoked. If no match
// is found and unmatchedHandler is non-nil, that handler is invoked. If no
// match is found and unmatchedHandler is nil, the message is dropped. If the
// hander returns a non-nil BusMessage, that message is sent as a reply to the
// message received.
func (b *Bus) HandleTypes(ctx context.Context, topic string, chanCap int,
	handlers map[int32]MessageHandler,
	unmatchedHandler MessageHandler,
) {
	in := make(chan *BusMessage, chanCap)
	b.Subscribe(topic, in)
	go func() {
		<-ctx.Done()
		b.Unsubscribe(topic, in)
		Drain(in)
	}()
	for msg := range in {
		handler := handlers[msg.Type]
		if handler == nil {
			handler = unmatchedHandler
		}
		if handler == nil {
			continue
		}
		if reply := handler(msg); reply != nil {
			b.SendReply(msg, reply)
		}
	}
}

// Drain drains all values from a channel, assuming that channel is closed.
func Drain[T any](c chan T) {
	for range c {
	}
}

// Error implements error for the Error proto
func (e *Error) Error() string {
	return e.GetDetail()
}

// DefaultReply creates a template reply by copying msg's topic and incrementing
// the message's type
func DefaultReply(msg *BusMessage) *BusMessage {
	return &BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
}
