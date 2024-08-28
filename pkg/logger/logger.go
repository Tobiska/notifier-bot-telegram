package logger

import (
	"context"
	"log/slog"
	"os"
)

type loggerSettings struct {
	handlerType int8
}

func NewLogger(opts ...option) *slog.Logger {
	logger := &loggerSettings{}
	for _, opt := range opts {
		opt(logger)
	}

	var handler slog.Handler
	if logger.handlerType == JsonHandler {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	return slog.New(handler)
}

func CtxWithSystemAttrs(ctx context.Context) context.Context {
	return WithAttrs(ctx, map[string]any{
		"pid": os.Getpid(),
	})
}
