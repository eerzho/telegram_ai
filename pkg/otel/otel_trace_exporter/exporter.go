package oteltraceexporter

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	ErrInvalidExporter = errors.New("invalid trace exporter type")
)

func MustNew(ctx context.Context, cfg Config) trace.SpanExporter {
	exp, err := New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return exp
}

func New(ctx context.Context, cfg Config) (trace.SpanExporter, error) {
	var traceExporter trace.SpanExporter
	var err error

	switch cfg.Exporter {
	case ExporterOTLPGRPC:
		traceExporter, err = otlptracegrpc.New(ctx)
	case ExporterOTLPHTTP:
		traceExporter, err = otlptracehttp.New(ctx)
	case ExporterConsole:
		traceExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	case ExporterNone:
		traceExporter = noopTraceExporter{}
	default:
		err = ErrInvalidExporter
	}

	if err != nil {
		return nil, fmt.Errorf("trace exporter: %w", err)
	}
	return traceExporter, nil
}

type noopTraceExporter struct{}

func (noopTraceExporter) ExportSpans(context.Context, []trace.ReadOnlySpan) error { return nil }
func (noopTraceExporter) Shutdown(context.Context) error                          { return nil }
