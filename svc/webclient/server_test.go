package webclient_test

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
)

type FakeFileEntry struct {
	ContentType string
	Data        []byte
}

type FakeFileServer map[string]FakeFileEntry

func (ffs FakeFileServer) RoundTrip(r *http.Request) (*http.Response, error) {
	entry, present := ffs[r.URL.String()]
	if !present {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     http.StatusText(http.StatusNotFound),
		}, nil
	}
	header := http.Header{}
	header.Set("Content-Type", entry.ContentType)
	header.Set("Content-Length", strconv.Itoa(len(entry.Data)))

	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Header:     header,
		Body:       io.NopCloser(bytes.NewBuffer(entry.Data)),
	}, nil
}

func NewFileServerClient(entries map[string]FakeFileEntry) *http.Client {
	return &http.Client{
		Transport: FakeFileServer(entries),
	}
}
