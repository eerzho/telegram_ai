package autootel

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

type noneTraceExporter struct{}

func (noneTraceExporter) ExportSpans(context.Context, []trace.ReadOnlySpan) error { return nil }
func (noneTraceExporter) Shutdown(context.Context) error                          { return nil }

func newTraceExporter(ctx context.Context) (trace.SpanExporter, error) {
	envProtocol := os.Getenv(ENVExporterOTLPProtocol)
	envExporter := os.Getenv(ENVTracesExporter)

	var exporter trace.SpanExporter
	var err error
	switch envExporter {
	case ExporterOTLP:
		switch envProtocol {
		case ProtocolHTTPProtobuff:
			exporter, err = otlptracehttp.New(ctx)
		case ProtocolHTTPJSON:
			exporter, err = otlptracehttp.New(ctx, otlptracehttp.WithHeaders(map[string]string{
				"Content-Type": "application/json",
			}))
		default:
			exporter, err = otlptracegrpc.New(ctx)
		}
	case ExporterConsole:
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	default:
		exporter = noneTraceExporter{}
	}

	if err != nil {
		return nil, fmt.Errorf("trace exporter: %w", err)
	}
	return exporter, nil
}

func newTracerProvider(res *resource.Resource, exp trace.SpanExporter) *trace.TracerProvider {
	processor := trace.NewBatchSpanProcessor(exp)
	provider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSpanProcessor(processor),
	)
	return provider
}
