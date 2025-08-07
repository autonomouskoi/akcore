package webclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/autonomouskoi/akcore"
	"github.com/autonomouskoi/akcore/modules/modutil"
	svc "github.com/autonomouskoi/akcore/svc/pb"
)

const InlineLimit = 64 * 1024

type Request struct {
	orig    *svc.WebclientHTTPRequest
	inPath  string
	outPath string
}

func NewRequest(pCtx *modutil.PluginContext, orig *svc.WebclientHTTPRequest) (*Request, error) {
	req := &Request{
		orig: orig,
	}

	if err := validateHTTPRequest(orig.Request); err != nil {
		return nil, fmt.Errorf("%w: invalid HTTP request: %w", akcore.ErrBadRequest, err)
	}

	var err error
	if disp, ok := orig.GetRequestBody().GetBodyAs().(*svc.BodyDisposition_Rwpath); ok {
		if req.inPath, err = translateRWPath(pCtx, disp); err != nil {
			return nil, fmt.Errorf("%w: request body path: %w", akcore.ErrBadRequest, err)
		}
	}
	if disp, ok := orig.GetResponseBody().GetBodyAs().(*svc.BodyDisposition_Rwpath); ok {
		if req.outPath, err = translateRWPath(pCtx, disp); err != nil {
			return nil, fmt.Errorf("%w: response body path: %w", akcore.ErrBadRequest, err)
		}
	}

	return req, nil
}

func validateHTTPRequest(req *svc.HTTPRequest) error {
	url, err := url.Parse(req.Url)
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err)
	}

	switch url.Scheme {
	case "http", "https":
		// cool
	default:
		return errors.New("invalid scheme")
	}

	switch req.GetMethod() {
	case http.MethodDelete, http.MethodGet, http.MethodHead, http.MethodOptions,
		http.MethodPatch, http.MethodPost, http.MethodPut:
		// cool
	default:
		return errors.New("invalid method")
	}

	// TODO: validate against per-plugin URL restrictions

	return nil
}

func translateRWPath(pCtx *modutil.PluginContext, disp *svc.BodyDisposition_Rwpath) (string, error) {
	if !strings.HasPrefix(disp.Rwpath, "/rwdata/") {
		return "", errors.New("not under /rwdata")
	}
	relPath := disp.Rwpath[len("/rwpath/"):]
	realPath := filepath.Join(pCtx.RWDataPath, relPath)
	if !strings.HasPrefix(realPath, pCtx.RWDataPath) {
		return "", fmt.Errorf("bad path: %s", disp.Rwpath)
	}
	return realPath, nil
}

func (req *Request) DoUsing(c *http.Client) (*svc.WebclientHTTPResponse, error) {
	reqBody, err := req.prepRequestBody()
	if err != nil {
		return nil, fmt.Errorf("%w: prepping request body: %w", akcore.ErrBadRequest, err)
	}
	defer reqBody.Close()
	r, err := http.NewRequest(req.orig.Request.Method, req.orig.Request.Url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("%w: creating HTTP request: %w", akcore.ErrBadRequest, err)
	}

	for k, v := range req.orig.Request.Header {
		r.Header[k] = v.GetValues()
	}
	r.Header.Set("User-Agent", "AutonomousKoi/"+akcore.Version)

	resp, err := c.Do(r)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	respP := &svc.WebclientHTTPResponse{
		Status:     resp.Status,
		StatusCode: int32(resp.StatusCode),
		Header:     map[string]*svc.StringValues{},
		Body:       &svc.BodyDisposition{},
	}
	for k, v := range resp.Header {
		sv := &svc.StringValues{}
		sv.Values = make([]string, len(v))
		copy(sv.Values, v)
		respP.Header[k] = sv
	}

	bodyDest, err := req.prepResponseBody(respP)
	if err != nil {
		return nil, fmt.Errorf("prepping response body: %w", err)
	}

	var respReader io.Reader = resp.Body
	if _, ok := req.orig.ResponseBody.BodyAs.(*svc.BodyDisposition_Inline); ok {
		respReader = io.LimitReader(respReader, InlineLimit+1) //+1 to detect going over
	}
	if _, err := io.Copy(bodyDest, respReader); err != nil {
		bodyDest.Close()
		return nil, fmt.Errorf("copying response body: %w", err)
	}

	if err := bodyDest.Close(); err != nil {
		return nil, fmt.Errorf("closing response body: %w", err)
	}

	return respP, nil
}

func (req *Request) prepRequestBody() (io.ReadCloser, error) {
	switch v := req.orig.GetRequestBody().GetBodyAs().(type) {
	case *svc.BodyDisposition_Inline:
		if len(v.Inline) > InlineLimit {
			return nil, fmt.Errorf("inline data exceeds %d bytes", InlineLimit)
		}
		r := bytes.NewReader(v.Inline)
		return io.NopCloser(r), nil
	case *svc.BodyDisposition_Rwpath:
		return os.Open(req.inPath)
	default:
		r := bytes.NewReader([]byte{})
		return io.NopCloser(r), nil
	}
}

func (req *Request) prepResponseBody(resp *svc.WebclientHTTPResponse) (io.WriteCloser, error) {
	switch v := req.orig.GetResponseBody().GetBodyAs().(type) {
	case *svc.BodyDisposition_Inline:
		il := &svc.BodyDisposition_Inline{}
		resp.Body.BodyAs = il
		buf := &bytes.Buffer{}
		return fnCloser{
			Writer: buf,
			closer: func() error {
				if buf.Len() > InlineLimit {
					return fmt.Errorf("inline data exceeds %d bytes", InlineLimit)
				}
				il.Inline = buf.Bytes()
				return nil
			},
		}, nil
	case *svc.BodyDisposition_Rwpath:
		bodyDest, err := os.Create(req.outPath)
		if err != nil {
			return nil, fmt.Errorf("creating %s: %w", v.Rwpath, err)
		}
		resp.Body.BodyAs = &svc.BodyDisposition_Rwpath{
			Rwpath: v.Rwpath,
		}
		return bodyDest, nil
	default:
		return fnCloser{
			Writer: io.Discard,
			closer: func() error { return nil },
		}, nil
	}
}
