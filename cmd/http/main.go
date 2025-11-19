package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eerzho/simpledi"
	"github.com/eerzho/telegram-ai/config"
	"github.com/eerzho/telegram-ai/internal/adapter/genkit"
	"github.com/eerzho/telegram-ai/internal/controller/http"
	health_check "github.com/eerzho/telegram-ai/internal/health/check"
	response_generate "github.com/eerzho/telegram-ai/internal/response/generate"
	summary_generate "github.com/eerzho/telegram-ai/internal/summary/generate"
	"github.com/eerzho/telegram-ai/pkg/httpserver"
	"github.com/eerzho/telegram-ai/pkg/logger"
	"github.com/go-playground/validator/v10"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	for _, definition := range definitions() {
		simpledi.Set(definition)
	}

	simpledi.Resolve()

	cfg := simpledi.Get[config.Config]("config")
	lgr := simpledi.Get[*slog.Logger]("logger")

	httpServer := httpserver.New(
		http.Handler(),
		cfg.HTTPServer,
	)

	serverErrs := make(chan error, 1)
	go func() {
		lgr.Info("starting http server", slog.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil {
			serverErrs <- err
		}
	}()

	select {
	case err := <-serverErrs:
		return fmt.Errorf("server: %w", err)
	case <-ctx.Done():
		lgr.Info("shutdown signal received")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	err := simpledi.Close()
	if err != nil {
		lgr.Warn("failed to close container", slog.Any("error", err))
	}

	lgr.Info("http server stopped")

	return nil
}

func definitions() []simpledi.Definition {
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
				return logger.MustNew(cfg.Logger)
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
			ID:   "healthCheckUsecase",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return health_check.NewUsecase(cfg.App.Version)
			},
		},
		{
			ID:   "responseGenerateUsecase",
			Deps: []string{"validate", "genkit"},
			New: func() any {
				validate := simpledi.Get[*validator.Validate]("validate")
				client := simpledi.Get[*genkit.Client]("genkit")
				return response_generate.NewUsecase(validate, client)
			},
		},
		{
			ID:   "summaryGenerateUsecase",
			Deps: []string{"validate", "genkit"},
			New: func() any {
				validate := simpledi.Get[*validator.Validate]("validate")
				client := simpledi.Get[*genkit.Client]("genkit")
				return summary_generate.NewUsecase(validate, client)
			},
		},
	}
}
