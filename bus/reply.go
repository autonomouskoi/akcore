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

func (b *Bus) WaitForReply(ctx context.Context, msg *BusMessage) *BusMessage {
	in := make(chan *BusMessage, 1)
	defer func() {
		for range in {
		}
	}() // drain channel
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
