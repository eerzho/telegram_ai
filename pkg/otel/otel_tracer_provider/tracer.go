package oteltracerprovider

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func MustNew(ctx context.Context, cfg Config, res *resource.Resource, exp trace.SpanExporter) *trace.TracerProvider {
	trc, err := New(ctx, cfg, res, exp)
	if err != nil {
		panic(err)
	}
	return trc
}

func New(ctx context.Context, cfg Config, res *resource.Resource, exp trace.SpanExporter) (*trace.TracerProvider, error) {
	provider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(exp,
			trace.WithBatchTimeout(cfg.BatchTimeout),
			trace.WithExportTimeout(cfg.ExportTimeout),
			trace.WithMaxExportBatchSize(cfg.MaxExportBatchSize),
			trace.WithMaxQueueSize(cfg.MaxQueueSize),
		),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return provider, nil
}
