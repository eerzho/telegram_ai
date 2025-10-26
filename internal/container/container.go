package container

import (
	"log/slog"

	"github.com/eerzho/simpledi"
	"github.com/eerzho/telegram-ai/config"
	"github.com/eerzho/telegram-ai/internal/adapter/genkit"
	"github.com/eerzho/telegram-ai/internal/generate_response"
	"github.com/eerzho/telegram-ai/internal/health_check"
	"github.com/eerzho/telegram-ai/pkg/logger"
	"github.com/go-playground/validator/v10"
)

func New() *simpledi.Container {
	c := simpledi.NewContainer()
	for _, def := range defs(c) {
		c.MustRegister(def)
	}
	c.MustResolve()
	return c
}

func defs(c *simpledi.Container) []simpledi.Def {
	return []simpledi.Def{
		{
			Key: "config",
			Ctor: func() any {
				return config.MustNew()
			},
		},
		{
			Key:  "logger",
			Deps: []string{"config"},
			Ctor: func() any {
				cfg := c.MustGet("config").(config.Config)
				return logger.MustNew(cfg.Logger)
			},
		},
		{
			Key: "validate",
			Ctor: func() any {
				return validator.New(validator.WithRequiredStructEnabled())
			},
		},
		{
			Key:  "genkit",
			Deps: []string{"config"},
			Ctor: func() any {
				cfg := c.MustGet("config").(config.Config)
				return genkit.New(cfg.Genkit)
			},
		},
		{
			Key:  "healthCheckUsecase",
			Deps: []string{"config"},
			Ctor: func() any {
				cfg := c.MustGet("config").(config.Config)
				return health_check.NewUsecase(cfg.App.Version)
			},
		},
		{
			Key:  "generateResponseUsecase",
			Deps: []string{"logger", "validate", "genkit"},
			Ctor: func() any {
				logger := c.MustGet("logger").(*slog.Logger)
				validate := c.MustGet("validate").(*validator.Validate)
				client := c.MustGet("genkit").(*genkit.Client)
				return generate_response.NewUsecase(logger, validate, client)
			},
		},
	}
}
