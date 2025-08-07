package webclient

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonomouskoi/akcore/modules/modutil"
	svc "github.com/autonomouskoi/akcore/svc/pb"
)

func TestValidateHTTPRequest(t *testing.T) {
	t.Parallel()

	validRequest := func() *svc.HTTPRequest {
		return &svc.HTTPRequest{
			Method: "GET",
			Url:    "https://autonomouskoi.org/path/",
			Header: map[string]*svc.StringValues{
				"Content-Type": {Values: []string{"text/plain"}},
			},
		}
	}

	t.Run("scheme", func(t *testing.T) {
		t.Parallel()

		for scheme, shouldError := range map[string]bool{
			"http": false, "https": false,
			"ssh": true, "ftp": true, "gopher": true,
		} {
			req := validRequest()
			req.Url = scheme + "://autonomouskoi.org/path/"
			err := validateHTTPRequest(req)
			if shouldError {
				require.ErrorContains(t, err, "invalid scheme", scheme)
			} else {
				require.NoError(t, err, scheme)
			}
		}
	})

	t.Run("method", func(t *testing.T) {
		t.Parallel()
		for method, shouldError := range map[string]bool{
			http.MethodConnect: true, http.MethodDelete: false, http.MethodGet: false,
			http.MethodHead: false, http.MethodOptions: false, http.MethodPatch: false,
			http.MethodPost: false, http.MethodPut: false, http.MethodTrace: true,
		} {
			req := validRequest()
			req.Method = method
			err := validateHTTPRequest(req)
			if shouldError {
				require.ErrorContains(t, err, "invalid method", method)
			} else {
				require.NoError(t, err, method)
			}
		}
	})

	t.Run("url parse", func(t *testing.T) {
		t.Parallel()
		for url, shouldError := range map[string]bool{
			"https://\\": true, "https://\n": true,
		} {
			req := validRequest()
			req.Url = url
			err := validateHTTPRequest(req)
			if shouldError {
				require.ErrorContains(t, err, "parsing URL", url)
			} else {
				require.NoError(t, err, url)
			}
		}
	})
}

func TestTranslateRWPath(t *testing.T) {
	t.Parallel()

	t.Run("not under", func(t *testing.T) {
		t.Parallel()

		for _, path := range []string{
			"/bloop",
			"/rwdata", "rwdata/", "/rwdatasomething",
		} {
			t.Run(path, func(t *testing.T) {
				_, err := translateRWPath(&modutil.PluginContext{}, &svc.BodyDisposition_Rwpath{
					Rwpath: path,
				})
				require.ErrorContains(t, err, "not under /rwdata")
			})
		}
	})

	t.Run("bad path", func(t *testing.T) {
		t.Parallel()

		for _, path := range []string{
			"/rwdata/.",
			"/rwdata/..",
			"/rwdata/../file",
			"/rwdata/something/../../file",
		} {
			t.Run(path, func(t *testing.T) {
				_, err := translateRWPath(&modutil.PluginContext{
					RWDataPath: "/plugin/data/",
				}, &svc.BodyDisposition_Rwpath{
					Rwpath: path,
				})
				require.ErrorContains(t, err, "bad path:")
			})
		}
	})

	t.Run("good path", func(t *testing.T) {
		for _, path := range []string{
			"/rwdata/file",
			"/rwdata/path/to/something",
		} {
			t.Run(path, func(t *testing.T) {
				_, err := translateRWPath(&modutil.PluginContext{
					RWDataPath: "/plugin/data/",
				}, &svc.BodyDisposition_Rwpath{
					Rwpath: path,
				})
				require.NoError(t, err)
			})
		}
	})
}

func TestPrepRequestBody(t *testing.T) {
	t.Parallel()

	testData := []byte("test data")

	t.Run("none", func(t *testing.T) {
		t.Parallel()
		req := &Request{
			orig: &svc.WebclientHTTPRequest{
				RequestBody: &svc.BodyDisposition{},
			},
		}
		r, err := req.prepRequestBody()
		require.NoError(t, err)
		require.NotNil(t, r)
		b, err := io.ReadAll(r)
		require.NoError(t, err, "reading")
		require.Empty(t, b)
	})

	t.Run("rwpath", func(t *testing.T) {
		dir := t.TempDir()
		inPath := filepath.Join(dir, "testfile")
		req := &Request{
			orig: &svc.WebclientHTTPRequest{
				RequestBody: &svc.BodyDisposition{
					BodyAs: &svc.BodyDisposition_Rwpath{
						Rwpath: inPath,
					},
				},
			},
			inPath: inPath,
		}
		require.NoError(t, os.WriteFile(inPath, testData, 0644), "writing test data")
		r, err := req.prepRequestBody()
		require.NoError(t, err)
		b, err := io.ReadAll(r)
		require.NoError(t, err, "reading")
		require.Equal(t, testData, b)

		// bad RWPath
		req.inPath = "/path/does/not/exist"
		_, err = req.prepRequestBody()
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("inline", func(t *testing.T) {
		req := &Request{
			orig: &svc.WebclientHTTPRequest{
				RequestBody: &svc.BodyDisposition{
					BodyAs: &svc.BodyDisposition_Inline{
						Inline: testData,
					},
				},
			},
		}
		r, err := req.prepRequestBody()
		require.NoError(t, err)
		b, err := io.ReadAll(r)
		require.NoError(t, err, "reading")
		require.Equal(t, testData, b)

		// too much data
		largeData := [InlineLimit + 1]byte{}
		req.orig.RequestBody.BodyAs = &svc.BodyDisposition_Inline{
			Inline: largeData[:],
		}
		_, err = req.prepRequestBody()
		require.ErrorContains(t, err, "inline data exceeds")
	})
}

func TestPrepResponseBody(t *testing.T) {
	t.Parallel()

	testData := []byte("test data")

	t.Run("none", func(t *testing.T) {
		req := &Request{
			orig: &svc.WebclientHTTPRequest{
				ResponseBody: &svc.BodyDisposition{},
			},
		}
		w, err := req.prepResponseBody(nil)
		require.NoError(t, err)
		require.NotNil(t, w)
		n, err := w.Write(testData)
		require.NoError(t, err, "writing")
		require.Len(t, testData, n)
		require.NoError(t, w.Close(), "closing")
	})

	t.Run("inline", func(t *testing.T) {
		req := &Request{
			orig: &svc.WebclientHTTPRequest{
				ResponseBody: &svc.BodyDisposition{
					BodyAs: &svc.BodyDisposition_Inline{},
				},
			},
		}
		resp := &svc.WebclientHTTPResponse{
			Body: &svc.BodyDisposition{},
		}
		w, err := req.prepResponseBody(resp)
		require.NoError(t, err)
		n, err := w.Write(testData)
		require.NoError(t, err, "writing")
		require.Len(t, testData, n)
		require.NoError(t, w.Close(), "closing")
		il, ok := resp.Body.BodyAs.(*svc.BodyDisposition_Inline)
		require.True(t, ok, "asserting body disposition inline")
		require.Equal(t, testData, il.Inline)

		// too much data
		largeData := [InlineLimit + 1]byte{}
		w, err = req.prepResponseBody(resp)
		require.NoError(t, err)
		_, err = w.Write(largeData[:])
		require.NoError(t, err, "writing")
		require.ErrorContains(t, w.Close(), "inline data exceeds")
	})

	t.Run("rwpath", func(t *testing.T) {
		dir := t.TempDir()
		outFile := filepath.Join(dir, "outfile")
		req := &Request{
			orig: &svc.WebclientHTTPRequest{
				ResponseBody: &svc.BodyDisposition{
					BodyAs: &svc.BodyDisposition_Rwpath{
						Rwpath: "/whatever",
					},
				},
			},
			outPath: outFile,
		}
		resp := &svc.WebclientHTTPResponse{
			Body: &svc.BodyDisposition{},
		}
		w, err := req.prepResponseBody(resp)
		require.NoError(t, err)
		require.NotNil(t, w)
		n, err := w.Write(testData)
		require.NoError(t, err, "writing")
		require.Len(t, testData, n)
		require.NoError(t, w.Close(), "closing")
		b, err := os.ReadFile(outFile)
		require.NoError(t, err, "reading file")
		require.Equal(t, testData, b)

		// bad path
		req.outPath = filepath.Join(dir, "does-not-exist", "valid-file")
		_, err = req.prepResponseBody(resp)
		require.ErrorIs(t, err, os.ErrNotExist)
	})
}
