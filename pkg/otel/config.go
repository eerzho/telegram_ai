package otel

import (
	otellogexporter "github.com/eerzho/telegram-ai/pkg/otel/otel_log_exporter"
	otelmeterprovider "github.com/eerzho/telegram-ai/pkg/otel/otel_meter_provider"
	otelmetricexporter "github.com/eerzho/telegram-ai/pkg/otel/otel_metric_exporter"
	otelresource "github.com/eerzho/telegram-ai/pkg/otel/otel_resource"
	oteltraceexporter "github.com/eerzho/telegram-ai/pkg/otel/otel_trace_exporter"
	oteltracerprovider "github.com/eerzho/telegram-ai/pkg/otel/otel_tracer_provider"
)

type Config struct {
	Resource otelresource.Config

	LogExporter otellogexporter.Config

	MetricExporter otelmetricexporter.Config
	MeterProvider  otelmeterprovider.Config

	TraceExporter  oteltraceexporter.Config
	TracerProvider oteltracerprovider.Config
}
