package modules

import (
	"net/http"
	"sync"
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
