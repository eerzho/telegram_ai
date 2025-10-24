package logger

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
)

var (
	ErrInvalidLogLevel = errors.New("invalid log level")
	ErrInvalidFormat   = errors.New("invalid format")
)

type Config struct {
	AppName    string `env:"APP_NAME,required"`
	AppVersion string `env:"APP_VERSION,required"`
	Level      string `env:"LOGGER_LEVEL"         envDefault:"info"` // debug, info, warn, error
	Format     string `env:"LOGGER_FORMAT"        envDefault:"json"` // text, json
}

func MustNew(cfg Config) *slog.Logger {
	l, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return l
}

func New(cfg Config) (*slog.Logger, error) {
	slogLevel, err := stringToSlogLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}
	handler, err := createHandler(cfg.Format, slogLevel, os.Stdout)
	if err != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}
	logger := slog.New(handler)
	logger = logger.With(
		slog.String("app_name", cfg.AppName),
		slog.String("app_version", cfg.AppVersion),
	)
	return logger, nil
}

func createHandler(format string, level slog.Level, w io.Writer) (slog.Handler, error) {
	opts := &slog.HandlerOptions{
		Level: level,
	}
	switch format {
	case "json":
		return slog.NewJSONHandler(w, opts), nil
	case "text":
		return slog.NewTextHandler(w, opts), nil
	default:
		return nil, ErrInvalidFormat
	}
}

func stringToSlogLevel(level string) (slog.Level, error) {
	switch level {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, ErrInvalidLogLevel
	}
}
