package modules

import (
	"net/http"
	"sync"

	"github.com/autonomouskoi/akcore/bus"
	"google.golang.org/protobuf/proto"
)

type handler struct {
	lock     sync.RWMutex
	handlers map[string]http.Handler
	mux      *http.ServeMux
}

func (mh *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mh.lock.RLock()
	mux := mh.mux
	mh.lock.RUnlock()
	mux.ServeHTTP(w, r)
}

func (mh *handler) Handle(path string, handler http.Handler) {
	mh.lock.Lock()
	defer mh.lock.Unlock()
	mh.mux.Handle(path, handler)
	mh.handlers[path] = handler
}

func (mh *handler) Remove(path string) {
	mh.lock.Lock()
	defer mh.lock.Unlock()
	if _, present := mh.handlers[path]; !present {
		return
	}
	delete(mh.handlers, path)
	mux := &http.ServeMux{}
	for path, handler := range mh.handlers {
		mux.Handle(path, handler)
	}
	mh.mux = mux
}

func webhooksHandler(b *bus.Bus, topic string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wcr := &bus.WebhookCallRequest{Params: map[string]*bus.WebhookValues{}}
		for k, v := range r.URL.Query() {
			wcr.Params[k] = &bus.WebhookValues{Values: v}
		}
		msg := &bus.BusMessage{
			Topic: topic,
			Type:  int32(bus.MessageTypeDirect_WEBHOOK_CALL_REQ),
		}
		msg.Message, _ = proto.Marshal(wcr)
		b.Send(msg)
	})
}
