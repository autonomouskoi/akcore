package akcore

type Logger interface {
	Error(string, ...any)
	Info(string, ...any)
	Debug(string, ...any)
}
