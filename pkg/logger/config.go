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
	FormatJSON = "json"
)

type Config struct {
	Level      LevelType         `env:"LOGGER_LEVEL"      envDefault:"info"`
	Format     FormatType        `env:"LOGGER_FORMAT"     envDefault:"json"`
	Attributes map[string]string `env:"LOGGER_ATTRIBUTES"                   envSeparator:","`
}

func (c Config) SlogLevel() slog.Level {
	switch c.Level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
