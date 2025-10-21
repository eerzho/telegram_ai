package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/eerzho/telegram-ai/pkg/httpserver"
	"github.com/eerzho/telegram-ai/pkg/logger"
	_ "github.com/joho/godotenv/autoload"
)

type App struct {
	Name    string `env:"APP_NAME"             envDefault:"Setting"`
	Version string `env:"APP_VERSION,required"`
}

type Config struct {
	App        App
	Logger     logger.Config
	HTTPServer httpserver.Config
}

func MustNew() Config {
	c, err := New()
	if err != nil {
		panic(err)
	}
	return c
}

func New() (Config, error) {
	const op = "config.New"

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("%s: %w", op, err)
	}

	return cfg, nil
}
