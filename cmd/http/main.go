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
	_ "github.com/eerzho/telegram-ai/docs"
	"github.com/eerzho/telegram-ai/internal/adapter/genkit"
	"github.com/eerzho/telegram-ai/internal/config"
	"github.com/eerzho/telegram-ai/internal/controller/http"
	healthcheck "github.com/eerzho/telegram-ai/internal/health/health_check"
	improvementgenerate "github.com/eerzho/telegram-ai/internal/improvement/improvement_generate"
	responsegenerate "github.com/eerzho/telegram-ai/internal/response/response_generate"
	summarygenerate "github.com/eerzho/telegram-ai/internal/summary/summary_generate"
	httpserver "github.com/eerzho/telegram-ai/pkg/http_server"
	"github.com/eerzho/telegram-ai/pkg/logger"
	otellogexporter "github.com/eerzho/telegram-ai/pkg/otel/otel_log_exporter"
	otelloggerprovider "github.com/eerzho/telegram-ai/pkg/otel/otel_logger_provider"
	otelmeterprovider "github.com/eerzho/telegram-ai/pkg/otel/otel_meter_provider"
	otelmetricexporter "github.com/eerzho/telegram-ai/pkg/otel/otel_metric_exporter"
	otelresource "github.com/eerzho/telegram-ai/pkg/otel/otel_resource"
	oteltraceexporter "github.com/eerzho/telegram-ai/pkg/otel/otel_trace_exporter"
	oteltracerprovider "github.com/eerzho/telegram-ai/pkg/otel/otel_tracer_provider"
	otelruntime "github.com/eerzho/telegram-ai/pkg/otel_help/otel_runtime"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/semaphore"
)

const (
	shutdownTimeout = 10 // second
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

	for _, definition := range definitions(ctx) {
		simpledi.Set(definition)
	}

	simpledi.Resolve()

	cfg := simpledi.Get[config.Config]("config")
	lgr := simpledi.Get[*slog.Logger]("logger")

	otelruntime.MustNew(cfg.OTELRuntime)
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
		lgr.InfoContext(ctx, "shutdown signal received")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	err := simpledi.Close()
	if err != nil {
		lgr.WarnContext(shutdownCtx, "failed to close container", slog.Any("error", err))
	}

	lgr.InfoContext(shutdownCtx, "http server stopped")

	return nil
}

func definitions(ctx context.Context) []simpledi.Definition {
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
			ID:   "otelResource",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return otelresource.MustNew(context.Background(), cfg.OTEL.Resource)
			},
		},
		{
			ID:   "otelLogExporter",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return otellogexporter.MustNew(ctx, cfg.OTEL.LogExporter)
			},
		},
		{
			ID:   "otelLoggerProvider",
			Deps: []string{"otelResource", "otelLogExporter"},
			New: func() any {
				resource := simpledi.Get[*resource.Resource]("otelResource")
				exporter := simpledi.Get[log.Exporter]("otelLogExporter")
				return otelloggerprovider.MustNew(ctx, resource, exporter)
			},
			Close: func() error {
				provider := simpledi.Get[*log.LoggerProvider]("otelLoggerProvider")
				return provider.Shutdown(ctx)
			},
		},
		{
			ID:   "otelMetricExporter",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return otelmetricexporter.MustNew(ctx, cfg.OTEL.MetricExporter)
			},
		},
		{
			ID:   "otelMeterProvider",
			Deps: []string{"config", "otelResource", "otelMetricExporter"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				resource := simpledi.Get[*resource.Resource]("otelResource")
				exporter := simpledi.Get[metric.Exporter]("otelMetricExporter")
				return otelmeterprovider.MustNew(ctx, cfg.OTEL.MeterProvider, resource, exporter)
			},
			Close: func() error {
				provider := simpledi.Get[*metric.MeterProvider]("otelMeterProvider")
				return provider.Shutdown(ctx)
			},
		},
		{
			ID:   "otelTraceExporter",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return oteltraceexporter.MustNew(ctx, cfg.OTEL.TraceExporter)
			},
		},
		{
			ID:   "otelTracerProvider",
			Deps: []string{"config", "otelResource", "otelTraceExporter"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				resource := simpledi.Get[*resource.Resource]("otelResource")
				exporter := simpledi.Get[trace.SpanExporter]("otelTraceExporter")
				return oteltracerprovider.MustNew(ctx, cfg.OTEL.TracerProvider, resource, exporter)
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
				return healthcheck.NewUsecase(cfg.App.Version)
			},
		},
		{
			ID:   "responseGenerateUsecase",
			Deps: []string{"generatorSem", "validate", "genkit"},
			New: func() any {
				generatorSem := simpledi.Get[*semaphore.Weighted]("generatorSem")
				validate := simpledi.Get[*validator.Validate]("validate")
				client := simpledi.Get[*genkit.Client]("genkit")
				return responsegenerate.NewUsecase(generatorSem, validate, client)
			},
		},
		{
			ID:   "summaryGenerateUsecase",
			Deps: []string{"generatorSem", "validate", "genkit"},
			New: func() any {
				generatorSem := simpledi.Get[*semaphore.Weighted]("generatorSem")
				validate := simpledi.Get[*validator.Validate]("validate")
				client := simpledi.Get[*genkit.Client]("genkit")
				return summarygenerate.NewUsecase(
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
				return improvementgenerate.NewUsecase(generatorSem, validate, client)
			},
		},
	}
}
