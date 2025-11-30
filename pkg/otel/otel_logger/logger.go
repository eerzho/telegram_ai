package otellogger

import (
	"context"

	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func MustNew(ctx context.Context, res *resource.Resource, exp log.Exporter) *log.LoggerProvider {
	lp, err := New(ctx, res, exp)
	if err != nil {
		panic(err)
	}
	return lp
}

func New(ctx context.Context, res *resource.Resource, exp log.Exporter) (*log.LoggerProvider, error) {
	processor := log.NewBatchProcessor(exp)

	provider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(processor),
	)

	return provider, nil
}
