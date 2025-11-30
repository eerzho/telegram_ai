package otelmetricexporter

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var (
	ErrInvalidExporter = errors.New("invalid metric exporter type")
)

func MustNew(ctx context.Context, cfg Config) metric.Exporter {
	exp, err := New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return exp
}

func New(ctx context.Context, cfg Config) (metric.Exporter, error) {
	var metricExporter metric.Exporter
	var err error

	switch cfg.Exporter {
	case ExporterOTLPGRPC:
		metricExporter, err = otlpmetricgrpc.New(ctx)
	case ExporterOTLPHTTP:
		metricExporter, err = otlpmetrichttp.New(ctx)
	case ExporterConsole:
		metricExporter, err = stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	case ExporterNone:
		metricExporter = noopMetricExporter{}
	default:
		err = ErrInvalidExporter
	}

	if err != nil {
		return nil, fmt.Errorf("metric exporter: %w", err)
	}
	return metricExporter, nil
}

type noopMetricExporter struct{}

func (noopMetricExporter) Temporality(metric.InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

func (noopMetricExporter) Aggregation(metric.InstrumentKind) metric.Aggregation {
	return nil
}

func (noopMetricExporter) Export(context.Context, *metricdata.ResourceMetrics) error {
	return nil
}

func (noopMetricExporter) ForceFlush(context.Context) error {
	return nil
}

func (noopMetricExporter) Shutdown(context.Context) error {
	return nil
}
