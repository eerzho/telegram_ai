package otelauto

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type noneMetricExporter struct{}

func (noneMetricExporter) Aggregation(metric.InstrumentKind) metric.Aggregation      { return nil }
func (noneMetricExporter) Export(context.Context, *metricdata.ResourceMetrics) error { return nil }
func (noneMetricExporter) ForceFlush(context.Context) error                          { return nil }
func (noneMetricExporter) Shutdown(context.Context) error                            { return nil }
func (noneMetricExporter) Temporality(metric.InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

func newMetricExporter(ctx context.Context) (metric.Exporter, error) {
	envProtocol := os.Getenv(ENVExporterOTLPProtocol)
	envExporter := os.Getenv(ENVMetricsExporter)

	var exporter metric.Exporter
	var err error
	switch envExporter {
	case ExporterOTLP:
		switch envProtocol {
		case ProtocolHTTPProtobuff:
			exporter, err = otlpmetrichttp.New(ctx)
		case ProtocolHTTPJSON:
			exporter, err = otlpmetrichttp.New(ctx, otlpmetrichttp.WithHeaders(map[string]string{
				"Content-Type": "application/json",
			}))
		default:
			exporter, err = otlpmetricgrpc.New(ctx)
		}
	case ExporterConsole:
		exporter, err = stdoutmetric.New(
			stdoutmetric.WithPrettyPrint(),
		)
	default:
		exporter = noneMetricExporter{}
	}

	if err != nil {
		return nil, fmt.Errorf("metric exporter: %w", err)
	}
	return exporter, nil
}

func newMeterProvider(res *resource.Resource, exp metric.Exporter) *metric.MeterProvider {
	reader := metric.NewPeriodicReader(exp)
	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	)
	return provider
}
