package logger

const (
	JsonHandler = iota
	TextHandler
)

type option func(l *loggerSettings)

func WithJSONHandler() func(l *loggerSettings) {
	return func(l *loggerSettings) {
		l.handlerType = JsonHandler
	}
}

func WithTextHandler() func(l *loggerSettings) {
	return func(l *loggerSettings) {
		l.handlerType = TextHandler
	}
}
