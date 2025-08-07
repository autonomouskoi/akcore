package webclient

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"

	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
	svc "github.com/autonomouskoi/akcore/svc/pb"
)

type fnCloser struct {
	io.Writer
	closer func() error
}

func (fc fnCloser) Close() error {
	return fc.closer()
}

type WebClient struct {
	http.Handler
	modutil.ModuleBase
	client       *http.Client
	cache        *cache
	cacheWebPath string
}

func New(deps *modutil.Deps, webPath string) (*WebClient, error) {
	wc := &WebClient{
		client:       deps.HttpClient,
		cacheWebPath: path.Join(webPath, "c") + "/",
	}
	wc.Log = deps.Log.NewForSource("svc.webclient")
	if wc.client == nil {
		wc.client = &http.Client{}
	}

	cacheDir := filepath.Join(deps.CachePath, "webclient")
	cache, err := newCache(cacheDir, wc.client)
	if err != nil {
		return nil, fmt.Errorf("creating cache: %w", err)
	}
	wc.cache = cache

	wc.Log.Info("initialized cache", "path", cacheDir, "web_path", wc.cacheWebPath)

	mux := http.NewServeMux()

	mux.Handle(wc.cacheWebPath, http.StripPrefix(wc.cacheWebPath, http.FileServer(http.Dir(cacheDir))))
	wc.Handler = mux

	return wc, nil
}

func (wc *WebClient) CloseModule(moduleID string) {

}

func (wc *WebClient) HandleRequestStaticDownload(msg *bus.BusMessage) *bus.BusMessage {
	reply := bus.DefaultReply(msg)
	cr := &svc.WebclientStaticDownloadRequest{}
	if reply.Error = wc.UnmarshalMessage(msg, cr); reply.Error != nil {
		return reply
	}
	filename, err := wc.cache.Get(cr.URL)
	if err != nil {
		reply.Error = &bus.Error{
			Detail: proto.String(err.Error()),
		}
		return reply
	}
	wc.MarshalMessage(reply, &svc.WebclientStaticDownloadResponse{
		Path: path.Join(wc.cacheWebPath, filename),
	})
	return reply
}

func (wc *WebClient) HandleRequest(pCtx modutil.PluginContext, msg *bus.BusMessage) *bus.BusMessage {
	reply := bus.DefaultReply(msg)
	hr := &svc.WebclientHTTPRequest{}
	if reply.Error = wc.UnmarshalMessage(msg, hr); reply.Error != nil {
		return reply
	}

	req, err := NewRequest(&pCtx, hr)
	if err != nil {
		reply.Error = bus.NewError(err)
		return reply
	}

	resp, err := req.DoUsing(wc.client)
	if err != nil {
		reply.Error = bus.NewError(err)
		return reply
	}

	wc.MarshalMessage(reply, resp)

	return reply
}
