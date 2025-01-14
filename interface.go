package akcore

// Logger is patterned after log/slog.Logger
type Logger interface {
	Error(string, ...any)
	Info(string, ...any)
	Debug(string, ...any)
}
