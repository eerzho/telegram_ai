package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/eerzho/setting/pkg/logger"
)

type App struct {
	Name    string `env:"APP_NAME" envDefault:"Setting"`
	Version string `env:"APP_VERSION,required"`
}

type Config struct {
	App    App
	Logger logger.Config
}

func Init() (Config, error) {
	const op = "config.Config.Init"

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("%s: %w", op, err)
	}

	return cfg, nil
}
