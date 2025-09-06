package modules

import (
	"net/http"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
	svc "github.com/autonomouskoi/akcore/svc/pb"
)

// handler wraps a ServeMux providing the ability to unhandle paths. It does
// this by building a new ServeMux that excludes the unregistered path
type handler struct {
	lock     sync.RWMutex
	handlers map[string]http.Handler
	mux      *http.ServeMux
}

// ServeHTTP serves a request using the built-in mux
func (mh *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mh.lock.RLock()
	mux := mh.mux
	mh.lock.RUnlock()
	mux.ServeHTTP(w, r)
}

// Handle registeres a handler with the mux
func (mh *handler) Handle(path string, handler http.Handler) {
	mh.lock.Lock()
	defer mh.lock.Unlock()
	mh.mux.Handle(path, handler)
	mh.handlers[path] = handler
}

// Remove a handler from the active mux
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

// create a handler that gets URL query parameters from from a web request and
// sends them on the topic as a WebhookCallRequest
func webhooksHandler(b *bus.Bus, topic string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wcr := &svc.WebhookCallEvent{Params: map[string]*svc.WebhookValues{}}
		for k, v := range r.URL.Query() {
			wcr.Params[k] = &svc.WebhookValues{Values: v}
		}
		msg := &bus.BusMessage{
			Topic: topic,
			Type:  int32(svc.MessageType_WEBHOOK_CALL_EVENT),
		}
		msg.Message, _ = proto.Marshal(wcr)
		b.Send(msg)
	})
}
