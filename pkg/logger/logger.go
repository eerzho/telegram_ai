package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
)

func New(cfg Config, handlers ...slog.Handler) *slog.Logger {
	slogLevel := cfg.SlogLevel()
	handler := createHandler(cfg.Format, slogLevel, os.Stdout)

	if len(handlers) > 0 {
		handlers = append(handlers, handler)
		handler = slogmulti.Fanout(handlers...)
	}

	logger := slog.New(handler)

	attrs := make([]any, 0, len(cfg.Attributes))
	for key, value := range cfg.Attributes {
		attrs = append(attrs, slog.String(key, value))
	}
	logger = logger.With(attrs...)

	return logger
}

func createHandler(format FormatType, level slog.Level, w io.Writer) slog.Handler {
	var handler slog.Handler
	switch format {
	case FormatText:
		handler = tint.NewHandler(w, &tint.Options{
			AddSource: true,
			Level:     level,
		})
	default:
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: level,
		})
	}
	return handler
}
