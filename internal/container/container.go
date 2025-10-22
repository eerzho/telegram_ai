package container

import (
	"github.com/eerzho/simpledi"
	"github.com/eerzho/telegram-ai/config"
	"github.com/eerzho/telegram-ai/internal/usecase"
	"github.com/eerzho/telegram-ai/pkg/logger"
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
			Key:  "healthUsecase",
			Deps: []string{"config"},
			Ctor: func() any {
				cfg := c.MustGet("config").(config.Config)
				return usecase.NewHealth(cfg.App.Version)
			},
		},
	}
}
