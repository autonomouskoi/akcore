package bus

import (
	"context"
	"math/rand"
	"time"
)

const (
	replyTimeout = time.Minute
)

type pendingReply struct {
	out     chan<- *BusMessage
	expires time.Time
}

// reply receives may be left dangling if no reply could be sent. Clean them up.
func (b *Bus) prunePendingReplies() {
	b.lock.Lock()
	defer b.lock.Unlock()
	now := time.Now()
	for id, pr := range b.pendingReplies {
		if pr.expires.After(now) {
			continue
		}
		close(pr.out)
		delete(b.pendingReplies, id)
	}
}

// periodically clean up up dangling reply handlers
func (b *Bus) prunePendingRepliesEvery(ctx context.Context, d time.Duration) {
	tick := time.NewTicker(d)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			b.prunePendingReplies()
		}
	}
}

// SendWithReply will send a message and wait for a reply via the provided
// channel. Replies are expired after one minute.
func (b *Bus) SendWithReply(msg *BusMessage, replyVia chan<- *BusMessage) {
	id := int64(rand.Int31())
	msg.ReplyTo = &id
	b.lock.Lock()
	b.pendingReplies[id] = pendingReply{
		expires: time.Now().Add(replyTimeout),
		out:     replyVia,
	}
	b.lock.Unlock()
	b.Send(msg)
}

// SendReply sends a message as a reply to another message. If no reply handler
// is waiting the message is dropped.
func (b *Bus) SendReply(from *BusMessage, msg *BusMessage) {
	b.lock.Lock()
	defer b.lock.Unlock()
	pr, present := b.pendingReplies[from.GetReplyTo()]
	if !present {
		return
	}
	msg.ReplyTo = from.ReplyTo
	pr.out <- msg
	close(pr.out)
	delete(b.pendingReplies, from.GetReplyTo())
}

// WaitForReply sends a message and waits for a reply, wrapping SendWithReply
func (b *Bus) WaitForReply(ctx context.Context, msg *BusMessage) *BusMessage {
	in := make(chan *BusMessage, 1)
	defer Drain(in)
	b.SendWithReply(msg, in)
	select {
	case <-ctx.Done():
		return &BusMessage{
			Error: &Error{
				Code: int32(CommonErrorCode_TIMEOUT),
			},
		}
	case reply := <-in:
		return reply
	}
}
