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
	oteltraceexporter "github.com/eerzho/telegram-ai/pkg/otel/otel_trace_exporter"
	otelmetricexporter "github.com/eerzho/telegram-ai/pkg/otel/otel_metric_exporter"
	otelmetricruntime "github.com/eerzho/telegram-ai/pkg/otel/otel_metric_runtime"
	otelresource "github.com/eerzho/telegram-ai/pkg/otel/otel_resource"
	oteltracer "github.com/eerzho/telegram-ai/pkg/otel/otel_tracer"
	otelmeter "github.com/eerzho/telegram-ai/pkg/otel/otel_meter"
	"github.com/go-playground/validator/v10"
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
			ID:   "generatorSem",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return semaphore.NewWeighted(cfg.App.GeneratorSemSize)
			},
		},
		{
			ID:   "otelResource",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return otelresource.MustNew(context.Background(), cfg.OTELResource)
			},
		},
		{
			ID:   "otelTraceExporter",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return oteltraceexporter.MustNew(context.Background(), cfg.OTELTraceExporter)
			},
		},
		{
			ID:   "otelTracer",
			Deps: []string{"config", "otelResource", "otelTraceExporter"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				resource := simpledi.Get[*resource.Resource]("otelResource")
				traceExporter := simpledi.Get[trace.SpanExporter]("otelTraceExporter")
				return oteltracer.MustNew(context.Background(), cfg.OTELTracer, resource, traceExporter)
			},
			Close: func() error {
				tracer := simpledi.Get[*trace.TracerProvider]("otelTracer")
				return tracer.Shutdown(context.Background())
			},
		},
		{
			ID:   "otelMetricExporter",
			Deps: []string{"config"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				return otelmetricexporter.MustNew(context.Background(), cfg.OTELMetricExporter)
			},
		},
		{
			ID:   "otelMeter",
			Deps: []string{"config", "otelResource", "otelMetricExporter"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				resource := simpledi.Get[*resource.Resource]("otelResource")
				metricExporter := simpledi.Get[metric.Exporter]("otelMetricExporter")
				return otelmeter.MustNew(context.Background(), cfg.OTELMeter, resource, metricExporter)
			},
			Close: func() error {
				meter := simpledi.Get[*metric.MeterProvider]("otelMeter")
				return meter.Shutdown(context.Background())
			},
		},
		{
			ID:   "otelMetricRuntime",
			Deps: []string{"otelMeter"},
			New: func() any {
				cfg := simpledi.Get[config.Config]("config")
				otelmetricruntime.MustNew(cfg.OTELMetricRuntime)
				return nil
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
