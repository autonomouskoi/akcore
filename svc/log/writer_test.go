package log_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonomouskoi/akcore/svc/log"
)

func TestWriter(t *testing.T) {
	t.Parallel()

	logDir := t.TempDir()

	// control time
	now := time.Date(2025, time.January, 1, 13, 25, 0, 0, time.Local)
	timeSrc := func() time.Time { return now }

	// attempt to create a writer with a bad dir
	neDir := filepath.Join(logDir, "does-not-exist")
	_, err := log.NewWriter(neDir, timeSrc)
	require.Error(t, err, "creating writer for non-existent dir")

	// create a writer successfully
	w, err := log.NewWriter(logDir, timeSrc)
	require.NoError(t, err, "creating logger")
	t.Cleanup(func() { w.Close() })

	// verify that the file was created
	stat, err := os.Stat(filepath.Join(logDir, "ak-20250101.log"))
	require.NoError(t, err, "statting log file")
	size := stat.Size()
	require.Zero(t, size, "log file should be empty")

	// write to the file
	written, err := w.Write([]byte{0, 1, 2, 3, 4})
	require.NoError(t, err, "writing")
	require.Equal(t, 5, written)

	// close it, create a new one
	require.NoError(t, w.Close(), "closing")
	w2, err := log.NewWriter(logDir, timeSrc)
	require.NoError(t, err, "creating logger")
	t.Cleanup(func() { w2.Close() })

	// write to it, verify that it's the same file
	written, err = w2.Write([]byte{5, 6, 7, 8, 9})
	require.NoError(t, err, "writing")
	require.Equal(t, 5, written)

	// verify that the data written was appended
	b, err := os.ReadFile(filepath.Join(logDir, "ak-20250101.log"))
	require.NoError(t, err, "reading file")
	require.Equal(t, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, b)

	// it's the next day, write again, verify that it's a new file
	now = time.Date(2025, time.January, 2, 13, 25, 0, 0, time.Local)
	written, err = w2.Write([]byte{10, 11, 12})
	require.NoError(t, err, "writing")
	require.Equal(t, 3, written)
	b, err = os.ReadFile(filepath.Join(logDir, "ak-20250102.log"))
	require.NoError(t, err, "reading file")
	require.Equal(t, []byte{10, 11, 12}, b)

	require.NoError(t, w2.Close())
}
