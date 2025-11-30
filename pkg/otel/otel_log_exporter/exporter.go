package otellogexporter

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
)

var (
	ErrInvalidExporter = errors.New("invalid log exporter type")
)

func MustNew(ctx context.Context, cfg Config) log.Exporter {
	exp, err := New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return exp
}

func New(ctx context.Context, cfg Config) (log.Exporter, error) {
	var logExporter log.Exporter
	var err error

	switch cfg.Exporter {
	case ExporterOTLPGRPC:
		logExporter, err = otlploggrpc.New(ctx)
	case ExporterOTLPHTTP:
		logExporter, err = otlploghttp.New(ctx)
	case ExporterConsole:
		logExporter, err = stdoutlog.New(stdoutlog.WithPrettyPrint())
	case ExporterNone:
		logExporter = noopLogExporter{}
	default:
		err = ErrInvalidExporter
	}

	if err != nil {
		return nil, fmt.Errorf("log exporter: %w", err)
	}
	return logExporter, nil
}

type noopLogExporter struct{}

func (noopLogExporter) Export(context.Context, []log.Record) error { return nil }
func (noopLogExporter) Shutdown(context.Context) error             { return nil }
func (noopLogExporter) ForceFlush(context.Context) error           { return nil }
