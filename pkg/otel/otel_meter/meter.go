package otelmeter

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func MustNew(ctx context.Context, cfg Config, res *resource.Resource, exp metric.Exporter) *metric.MeterProvider {
	mp, err := New(ctx, cfg, res, exp)
	if err != nil {
		panic(err)
	}
	return mp
}

func New(ctx context.Context, cfg Config, res *resource.Resource, exp metric.Exporter) (*metric.MeterProvider, error) {
	reader := metric.NewPeriodicReader(exp)

	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
		metric.WithCardinalityLimit(cfg.CardinalityLimit),
	)

	otel.SetMeterProvider(provider)

	return provider, nil
}
