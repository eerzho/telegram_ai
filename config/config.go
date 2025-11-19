package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/eerzho/telegram-ai/internal/adapter/genkit"
	"github.com/eerzho/telegram-ai/internal/adapter/postgres"
	"github.com/eerzho/telegram-ai/internal/adapter/valkey"
	"github.com/eerzho/telegram-ai/pkg/cors"
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
	CORS       cors.Config
	Genkit     genkit.Config
	Valkey     valkey.Config
	Postgres   postgres.Config
}

func MustNew() Config {
	c, err := New()
	if err != nil {
		panic(err)
	}
	return c
}

func New() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}
