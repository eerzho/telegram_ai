package logger

import (
	"context"
	"log/slog"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

type ContextHandler struct {
	handler slog.Handler
}

func NewContextHandler(handler slog.Handler) *ContextHandler {
	return &ContextHandler{handler: handler}
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ContextHandler) Handle(ctx context.Context, record slog.Record) error {
	if id, ok := ctx.Value(RequestIDKey).(string); ok && id != "" {
		record.AddAttrs(slog.String("request_id", id))
	}
	return h.handler.Handle(ctx, record)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewContextHandler(h.handler.WithAttrs(attrs))
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return NewContextHandler(h.handler.WithGroup(name))
}
