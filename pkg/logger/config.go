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
	Level      LevelType  `env:"LOGGER_LEVEL"         envDefault:"info"`
	Format     FormatType `env:"LOGGER_FORMAT"        envDefault:"json"`
}

func (c Config) SlogLevel() (slog.Level, error) {
	switch c.Level {
	case LevelDebug:
		return slog.LevelDebug, nil
	case LevelInfo:
		return slog.LevelInfo, nil
	case LevelWarn:
		return slog.LevelWarn, nil
	case LevelError:
		return slog.LevelError, nil
	default:
		return 0, ErrInvalidLogLevel
	}
}
