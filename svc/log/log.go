package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	svc "github.com/autonomouskoi/akcore/svc/pb"
)

const (
	// DefaultLogLevel is the default logging level
	DefaultLogLevel = slog.LevelInfo
)

// Logger has logging functions in the style of log/slog.Logger
type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
}

// MasterLogger provides loggers
type MasterLogger struct {
	w       io.WriteCloser
	level   slog.Level
	logsDir string
}

// New creates a new MasterLogger
func New(logsDir string, cfg *svc.Config) (*MasterLogger, error) {
	w, err := NewWriter(logsDir, time.Now)
	if err != nil {
		return nil, fmt.Errorf("creating writer: %w", err)
	}
	level := DefaultLogLevel
	if cfg.LogLevel != nil {
		switch *cfg.LogLevel {
		case svc.LogLevel_DEBUG:
			level = slog.LevelDebug
		case svc.LogLevel_INFO:
			level = slog.LevelInfo
		case svc.LogLevel_WARN:
			level = slog.LevelWarn
		case svc.LogLevel_ERROR:
			level = slog.LevelError
		}
	}
	return &MasterLogger{
		w:       w,
		level:   level,
		logsDir: logsDir,
	}, nil
}

// Close the underlying writer
func (ml *MasterLogger) Close() error {
	ml.cleanup()
	return ml.w.Close()
}

// NewForSource provides a logger for a given source. That source may have a
// distinct configuration.
func (ml *MasterLogger) NewForSource(src string) Logger {
	level := ml.level
	logger := slog.New(slog.NewTextHandler(
		ml.w,
		&slog.HandlerOptions{
			Level: level,
		},
	)).With("source", src)

	return logger
}

func (ml *MasterLogger) cleanup() {
	logger := ml.NewForSource("master_logger.cleanup")
	entries, err := os.ReadDir(ml.logsDir)
	if err != nil {
		logger.Error("listing log files", "error", err.Error())
		return
	}
	cutoff := time.Now().Add(-time.Hour * 24 * 30)
	logger.Debug("deleting old log files", "cutoff", cutoff)
	for _, entry := range entries {
		fullpath := filepath.Join(ml.logsDir, entry.Name())
		if entry.IsDir() {
			logger.Debug("skipping dir", "dir", fullpath)
			continue
		}
		stat, err := os.Stat(fullpath)
		if err != nil {
			logger.Error("statting file", "path", fullpath, "error", err.Error())
			continue
		}
		if stat.ModTime().Before(cutoff) {
			logger.Debug("deleting log file", "path", fullpath)
			if err := os.Remove(fullpath); err != nil {
				logger.Error("deleting olg file", "path", fullpath, "error", err.Error())
				continue
			}
		} else {
			logger.Debug("keeping file", "path", fullpath)
		}
	}
}
