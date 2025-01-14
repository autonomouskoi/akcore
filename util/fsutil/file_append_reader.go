package fsutil

import (
	"fmt"
	"io"
	"os"
	"time"
)

// PollingFileAppendReader polls a file for updates. When the size of the file
// changes it reads the new content and makes it available via Read()
type PollingFileAppendReader struct {
	io.Reader
	w        *io.PipeWriter
	path     string
	interval time.Duration
	lastSize int64
}

// NewPollingFileAppendReader creates a new reader for the provided path which
// polls at the desired interval
func NewPollingFileAppendReader(path string, interval time.Duration) (*PollingFileAppendReader, error) {
	// test open
	infh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	infh.Close()
	pipeReader, pipeWriter := io.Pipe()
	r := &PollingFileAppendReader{
		Reader:   pipeReader,
		w:        pipeWriter,
		path:     path,
		interval: interval,
	}
	go r.pollLoop()
	return r, nil
}

// Close the reader
func (r *PollingFileAppendReader) Close() error {
	r.closeWithError(nil)
	return nil
}

func (r *PollingFileAppendReader) closeWithError(err error) {
	r.w.CloseWithError(err)
	r.w = nil
}

func (r *PollingFileAppendReader) pollLoop() {
	for {
		if err := r.poll(); err != nil {
			r.closeWithError(err)
		}
		if r.w == nil {
			break
		}
		time.Sleep(r.interval)
	}
}

func (r *PollingFileAppendReader) poll() error {
	fi, err := os.Stat(r.path)
	if err != nil {
		return err
	}
	currentSize := fi.Size()
	if currentSize <= r.lastSize {
		return nil
	}
	infh, err := os.Open(r.path)
	if err != nil {
		return err
	}
	defer infh.Close()
	if _, err := infh.Seek(r.lastSize, 0); err != nil {
		return fmt.Errorf("seeking: %w", err)
	}
	n, err := io.Copy(r.w, infh)
	if err != nil {
		return err
	}
	r.lastSize += n
	return nil
}
