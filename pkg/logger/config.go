package logger

import "log/slog"

type LevelType string

const (
	LevelDebug LevelType = "debug"
	LevelInfo  LevelType = "info"
	LevelWarn  LevelType = "warn"
	LevelError LevelType = "error"
)

type FormatType string

const (
	FormatText = "text"
	FormatJson = "json"
)

type Config struct {
	AppName    string     `env:"APP_NAME,required"`
	AppVersion string     `env:"APP_VERSION,required"`
	Level      LevelType  `env:"LOGGER_LEVEL"         envDefault:"info"` // debug, info, warn, error
	Format     FormatType `env:"LOGGER_FORMAT"        envDefault:"json"` // text, json
}

func (c Config) SlogLevel() (slog.Level, error) {
	switch c.Level {
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
