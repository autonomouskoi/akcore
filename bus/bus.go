package bus

import (
	"context"
	"sync"
	"time"

	"github.com/autonomouskoi/datastruct/mapset"
)

type Bus struct {
	lock           sync.Mutex
	subs           map[string]mapset.MapSet[chan<- *BusMessage]
	pendingReplies map[int64]pendingReply
}

func New(ctx context.Context) *Bus {
	b := &Bus{
		subs:           map[string]mapset.MapSet[chan<- *BusMessage]{},
		pendingReplies: map[int64]pendingReply{},
	}

	go b.prunePendingRepliesEvery(ctx, time.Minute)
	go func() {
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

func (b *Bus) HasTopic(topic string) bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return len(b.subs[topic]) > 0
}

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

func (b *Bus) Send(msg *BusMessage) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for recv := range b.subs[msg.GetTopic()] {
		select {
		case recv <- msg:
			// cool
		default:
			// discarding, also cool
		}
	}
}

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

type MessageHandler func(*BusMessage) *BusMessage

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

func Drain[T any](c chan T) {
	for range c {
	}
}

func (e *Error) Error() string {
	return e.GetDetail()
}
