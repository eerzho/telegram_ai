package otelauto

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

type noneLogExporter struct{}

func (noneLogExporter) Export(context.Context, []log.Record) error { return nil }
func (noneLogExporter) Shutdown(context.Context) error             { return nil }
func (noneLogExporter) ForceFlush(context.Context) error           { return nil }

func newLogExporter(ctx context.Context) (log.Exporter, error) {
	envProtocol := os.Getenv(ENVExporterOTLPProtocol)
	envExporter := os.Getenv(ENVLogsExporter)

	var exporter log.Exporter
	var err error
	switch envExporter {
	case ExporterOTLP:
		switch envProtocol {
		case ProtocolHTTPProtobuff:
			exporter, err = otlploghttp.New(ctx)
		case ProtocolHTTPJSON:
			exporter, err = otlploghttp.New(ctx, otlploghttp.WithHeaders(map[string]string{
				"Content-Type": "application/json",
			}))
		default:
			exporter, err = otlploggrpc.New(ctx)
		}
	case ExporterConsole:
		exporter, err = stdoutlog.New(
			stdoutlog.WithPrettyPrint(),
		)
	default:
		exporter = noneLogExporter{}
	}

	if err != nil {
		return nil, fmt.Errorf("log exporter: %w", err)
	}
	return exporter, nil
}

func newLoggerProvider(res *resource.Resource, exp log.Exporter) *log.LoggerProvider {
	processor := log.NewBatchProcessor(exp)
	provider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(processor),
	)
	return provider
}
