package log_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/autonomouskoi/akcore/svc/log"
	svc "github.com/autonomouskoi/akcore/svc/pb"
	"github.com/stretchr/testify/require"
)

func TestMasterLogger(t *testing.T) {
	t.Parallel()

	logDir := t.TempDir()

	// create our master logger
	ml, err := log.New(logDir, &svc.Config{})
	require.NoError(t, err, "creating master logger")
	defer ml.Close()

	// there should be an empty file in there. Discover the name
	entries, err := os.ReadDir(logDir)
	require.NoError(t, err, "reading log dir")
	require.Len(t, entries, 1)
	entry := entries[0]
	require.False(t, entry.IsDir())
	require.True(t, strings.HasSuffix(entry.Name(), ".log"))
	logPath := filepath.Join(logDir, entry.Name())
	t.Logf("log path: %q", logPath)

	getSize := func() int64 {
		t.Helper()
		stat, err := os.Stat(logPath)
		require.NoError(t, err, "statting log file")
		return stat.Size()
	}

	// it should be empty at first
	require.Zero(t, getSize())

	// get a logger and write to it, there should be some data
	logger := ml.NewForSource(t.Name())
	logger.Info("info12345")
	size := getSize()
	require.NotZero(t, size)

	b, err := os.ReadFile(logPath)
	require.NoError(t, err, "reading log file")
	require.Contains(t, string(b), "info12345")

	// A debug message should not get emitted
	logger.Debug("debug12345")
	require.Equal(t, size, getSize())

	// An error should be
	logger.Error("error12345")
	require.Greater(t, getSize(), size)
	size = getSize()

	// close the logger, create one with a different log level
	require.NoError(t, ml.Close(), "closing ml")
	ml, err = log.New(logDir, &svc.Config{LogLevel: svc.LogLevel_ERROR.Enum()})
	require.NoError(t, err, "creating logger")

	// a message at a lower level is not emitted
	logger = ml.NewForSource(t.Name())
	logger.Info("info67890")
	require.Equal(t, size, getSize())

	// but an error is
	logger.Error("error67890")
	require.Greater(t, getSize(), size)
}
