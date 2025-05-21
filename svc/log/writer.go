package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// NewWriter creates an io.WriteCloser that will write to a file in the specified
// directory in a file named after the current date. The date is determined by
// timeSource. If timeSource is nil, time.Now is used.
func NewWriter(dir string, timeSource func() time.Time) (io.WriteCloser, error) {
	if timeSource == nil {
		timeSource = time.Now
	}
	w := &writer{
		dir: dir,
		day: -1,
		now: timeSource,
	}
	return w, w.open(w.now())
}

type writer struct {
	lock  sync.Mutex
	outfh *os.File
	dir   string
	day   int
	now   func() time.Time
}

// Write like io.Writer. Calls to Write are concurrency-safe
func (w *writer) Write(b []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	now := w.now()

	if day := now.Day(); day != w.day {
		if err := w.close(); err != nil {
			return -1, fmt.Errorf("closing log: %w", err)
		}
		if err := w.open(now); err != nil {
			return -1, fmt.Errorf("opening: %w", err)
		}
	}
	return w.outfh.Write(b)
}

func (w *writer) open(now time.Time) error {
	logPath := filepath.Join(w.dir, now.Format("ak-20060102.log"))
	outfh, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		return err
	}
	w.outfh = outfh
	w.day = now.Day()
	return nil
}

// Close the currently open log file. It is safe to close if no file is open
func (w *writer) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.close()
}

func (w *writer) close() error {
	if w.outfh == nil {
		return nil
	}
	outfh := w.outfh
	w.outfh = nil
	if err := outfh.Sync(); err != nil {
		return fmt.Errorf("syncing file: %w", err)
	}
	return outfh.Close()
}
