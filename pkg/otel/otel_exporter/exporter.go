package otelexporter

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
	ErrInvalidExporter = errors.New("invalid exporter type")
)

func MustNew(ctx context.Context, cfg Config) trace.SpanExporter {
	exp, err := New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return exp
}

func New(ctx context.Context, cfg Config) (trace.SpanExporter, error) {
	var exporter trace.SpanExporter
	var err error

	switch cfg.Exporter {
	case ExporterOTLPGRPC:
		exporter, err = otlptracegrpc.New(ctx)
	case ExporterOTLPHTTP:
		exporter, err = otlptracehttp.New(ctx)
	case ExporterConsole:
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	case ExporterNone:
		exporter = noopExporter{}
	default:
		err = ErrInvalidExporter
	}

	if err != nil {
		return nil, fmt.Errorf("exporter: %w", err)
	}
	return exporter, nil
}

type noopExporter struct{}

func (noopExporter) ExportSpans(context.Context, []trace.ReadOnlySpan) error { return nil }
func (noopExporter) Shutdown(context.Context) error                          { return nil }
