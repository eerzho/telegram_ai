package otelexporter

type ExporterType string

const (
	ExporterOTLPGRPC ExporterType = "otlp-grpc"
	ExporterOTLPHTTP ExporterType = "otlp-http"
	ExporterConsole  ExporterType = "console"
	ExporterNone     ExporterType = "none"
)

type Config struct {
	Exporter ExporterType `env:"OTEL_TRACES_EXPORTER" envDefault:"otlp-grpc"`
}
