package container

import (
	autootel "github.com/eerzho/goiler/pkg/auto_otel"
	"github.com/eerzho/goiler/pkg/logger"
	"github.com/eerzho/simpledi"
	"github.com/eerzho/telegram_ai/internal/adapter/genkit"
	"github.com/eerzho/telegram_ai/internal/config"
	generateimprovement "github.com/eerzho/telegram_ai/internal/improvement/generate_improvement"
	healthcheck "github.com/eerzho/telegram_ai/internal/monitoring/health_check"
	generateresponse "github.com/eerzho/telegram_ai/internal/response/generate_response"
	generatesummary "github.com/eerzho/telegram_ai/internal/summary/generate_summary"
	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/semaphore"
)

func Definitions() []simpledi.Definition {
	return []simpledi.Definition{
		{
			ID: "config",
			New: func() any {
				return config.MustNew()
			},
		},
		{
			ID:   "logger",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return logger.New(cfg.Logger, autootel.NewSlogHandler())
			},
		},
		{
			ID: "validate",
			New: func() any {
				return validator.New(validator.WithRequiredStructEnabled())
			},
		},
		{
			ID:   "genkit",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return genkit.New(cfg.Genkit)
			},
		},
		{
			ID:   "generatorSem",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return semaphore.NewWeighted(cfg.App.GeneratorSemSize)
			},
		},
		{
			ID:   "healthCheckUsecase",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return healthcheck.NewUsecase(cfg.App)
			},
		},
		{
			ID:   "responseGenerateUsecase",
			Deps: []string{"generatorSem", "validate", "genkit"},
			New: func() any {
				generatorSem := simpledi.Get[*semaphore.Weighted]("generatorSem")
				validate := simpledi.Get[*validator.Validate]("validate")
				client := simpledi.Get[*genkit.Client]("genkit")
				return generateresponse.NewUsecase(generatorSem, validate, client)
			},
		},
		{
			ID:   "summaryGenerateUsecase",
			Deps: []string{"generatorSem", "validate", "genkit"},
			New: func() any {
				generatorSem := simpledi.Get[*semaphore.Weighted]("generatorSem")
				validate := simpledi.Get[*validator.Validate]("validate")
				client := simpledi.Get[*genkit.Client]("genkit")
				return generatesummary.NewUsecase(
					generatorSem,
					validate,
					client,
				)
			},
		},
		{
			ID:   "improvementGenerateUsecase",
			Deps: []string{"generatorSem", "validate", "genkit"},
			New: func() any {
				generatorSem := simpledi.Get[*semaphore.Weighted]("generatorSem")
				validate := simpledi.Get[*validator.Validate]("validate")
				client := simpledi.Get[*genkit.Client]("genkit")
				return generateimprovement.NewUsecase(generatorSem, validate, client)
			},
		},
	}
}
