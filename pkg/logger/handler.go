package logger

import (
	"context"
	"log/slog"
)

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

func WithAttr(ctx context.Context, key string, value any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	attr := slog.Any(key, value)

	if v, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(ctx, slogFields, v)
	}

	var v []slog.Attr
	v = append(v, attr)
	return context.WithValue(ctx, slogFields, v)
}

func WithAttrs(updatableCtx context.Context, appendAttrs map[string]any) context.Context {
	if updatableCtx == nil {
		updatableCtx = context.Background()
	}

	var attrs []slog.Attr
	if v, ok := updatableCtx.Value(slogFields).([]slog.Attr); ok {
		attrs = v
	}

	for key, value := range appendAttrs {
		attr := slog.Any(key, value)
		attrs = append(attrs, attr)
	}

	for _, attr := range attrs {
		updatableCtx = context.WithValue(updatableCtx, slogFields, attr)
	}
	return updatableCtx
}
