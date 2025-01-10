package web

import (
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/autonomouskoi/akcore"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/akcore/web/webutil"
)

const (
	EnvLocalContentPath = "AK_CONTENT_WEB"
)

type Web struct {
	http.Handler
	basePattern string
	mux         *http.ServeMux
	log         akcore.Logger
}

//go:embed web.zip
var webZip []byte

func New(basePattern string, deps *modutil.Deps) *Web {
	log := deps.Log.With("module", "web")
	mux := http.NewServeMux()

	mux.Handle("/ws", newWS(deps))
	mux.HandleFunc("/build.json", handleBuildJSON)

	fontsPath := filepath.Join(deps.StoragePath, "fonts")
	if _, err := os.Stat(fontsPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(fontsPath, 0755); err != nil {
				log.Error("creating fonts dir", "path", fontsPath, "error", err.Error())
			}
		} else {
			log.Error("checking fonts path", "path", fontsPath, "error", err.Error())
		}
	}
	mux.Handle("/fonts/", http.StripPrefix("/fonts", http.FileServer(http.Dir(fontsPath))))

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

func handleBuildJSON(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(map[string]string{
		"Software": "AutonomousKoi",
		"Build":    "v" + akcore.Version,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (w *Web) Handle(path string, handler http.Handler) {
	w.mux.Handle(path, handler)
}
