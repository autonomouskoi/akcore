package web

import (
	_ "embed"
	"net/http"

	"github.com/autonomouskoi/akcore"
	"github.com/autonomouskoi/akcore/config"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/akcore/web/webutil"
)

const (
	EnvLocalContentPath = "AK_WEB_CONTENT"
)

type Web struct {
	http.Handler
	basePattern string
	mux         *http.ServeMux
	log         akcore.Logger
}

//go:embed web.zip
var webZip []byte

func New(cfg *config.Web, basePattern string, deps *modutil.Deps) *Web {
	log := deps.Log.With("module", "web")
	mux := http.NewServeMux()

	mux.Handle("/ws", newWS(deps))

	fs, err := webutil.ZipOrEnvPath(EnvLocalContentPath, webZip)
	if err != nil {
		panic("CRAP")
	}
	mux.Handle("/", http.FileServer(fs))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//_, pattern := mux.Handler(r)
		//log.Debug("HTTP Request", "method", r.Method, "path", r.URL.Path, "pattern", pattern)
		mux.ServeHTTP(w, r)
	})

	return &Web{
		Handler:     handler,
		basePattern: basePattern,
		mux:         mux,
		log:         log,
	}
}

func (w *Web) Handle(path string, handler http.Handler) {
	w.mux.Handle(path, handler)
}
