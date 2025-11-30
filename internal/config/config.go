package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/eerzho/telegram-ai/internal/adapter/genkit"
	bodysize "github.com/eerzho/telegram-ai/pkg/body_size"
	"github.com/eerzho/telegram-ai/pkg/cors"
	httpserver "github.com/eerzho/telegram-ai/pkg/http_server"
	"github.com/eerzho/telegram-ai/pkg/logger"
	otelexporter "github.com/eerzho/telegram-ai/pkg/otel/otel_exporter"
	otelresource "github.com/eerzho/telegram-ai/pkg/otel/otel_resource"
	oteltracer "github.com/eerzho/telegram-ai/pkg/otel/otel_tracer"
	"github.com/joho/godotenv"
)

type App struct {
	Name             string `env:"APP_NAME,required"`
	Version          string `env:"APP_VERSION,required"`
	GeneratorSemSize int64  `env:"APP_GENERATOR_SEM_SIZE" envDefault:"1000"`
}

type Config struct {
	App          App
	Logger       logger.Config
	HTTPServer   httpserver.Config
	CORS         cors.Config
	Genkit       genkit.Config
	BodySize     bodysize.Config
	OTELResource otelresource.Config
	OTELExporter otelexporter.Config
	OTELTracer   oteltracer.Config
}

func MustNew() Config {
	c, err := New()
	if err != nil {
		panic(err)
	}
	return c
}

func New() (Config, error) {
	_ = godotenv.Load()

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}
