package otelauto

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	ENVServiceName          string = "OTEL_SERVICE_NAME"
	ENVGoDetectors          string = "OTEL_GO_DETECTORS"
	ENVResourceSchemaURL    string = "OTEL_RESOURCE_SCHEMA_URL"
	ENVTracesExporter       string = "OTEL_TRACES_EXPORTER"
	ENVLogsExporter         string = "OTEL_LOGS_EXPORTER"
	ENVMetricsExporter      string = "OTEL_METRICS_EXPORTER"
	ENVExporterOTLPProtocol string = "OTEL_EXPORTER_OTLP_PROTOCOL"

	ExporterOTLP    string = "otlp"
	ExporterConsole string = "console"
	ExporterNone    string = "none"

	ProtocolGRPC          string = "grpc"
	ProtocolHTTPProtobuff string = "http/protobuf"
	ProtocolHTTPJSON      string = "http/json"

	DetectorAll                       string = "all"
	DetectorENV                       string = "env"
	DetectorHost                      string = "host"
	DetectorHostID                    string = "host_id"
	DetectorTelemetrySDK              string = "telemetry_sdk"
	DetectorSchemaURL                 string = "schema_url"
	DetectorOS                        string = "os"
	DetectorOSType                    string = "os_type"
	DetectorOSDescription             string = "os_description"
	DetectorProcess                   string = "process"
	DetectorProcessPID                string = "process_pid"
	DetectorProcessExecutableName     string = "process_executable_name"
	DetectorProcessExecutablePath     string = "process_executable_path"
	DetectorProcessCommandArgs        string = "process_command_args"
	DetectorProcessOwner              string = "process_owner"
	DetectorProcessRuntimeName        string = "process_runtime_name"
	DetectorProcessRuntimeVersion     string = "process_runtime_version"
	DetectorProcessRuntimeDescription string = "process_runtime_description"
	DetectorContainer                 string = "container"
	DetectorContainerID               string = "container_id"
)

func newResource(ctx context.Context) (*resource.Resource, error) {
	envDetectors := os.Getenv(ENVGoDetectors)
	if envDetectors == "" {
		return newResourceAll(ctx)
	}
	detectors := strings.Split(envDetectors, ",")
	opts := make([]resource.Option, 0, len(detectors))
	for _, detector := range detectors {
		switch strings.TrimSpace(detector) {
		case DetectorAll:
			return newResourceAll(ctx)
		case DetectorENV:
			opts = append(opts, resource.WithFromEnv())
		case DetectorHost:
			opts = append(opts, resource.WithHost())
		case DetectorHostID:
			opts = append(opts, resource.WithHostID())
		case DetectorTelemetrySDK:
			opts = append(opts, resource.WithTelemetrySDK())
		case DetectorSchemaURL:
			opts = append(opts, resource.WithSchemaURL(os.Getenv(string(ENVResourceSchemaURL))))
		case DetectorOS:
			opts = append(opts, resource.WithOS())
		case DetectorOSType:
			opts = append(opts, resource.WithOSType())
		case DetectorOSDescription:
			opts = append(opts, resource.WithOSDescription())
		case DetectorProcess:
			opts = append(opts, resource.WithProcess())
		case DetectorProcessPID:
			opts = append(opts, resource.WithProcessPID())
		case DetectorProcessExecutableName:
			opts = append(opts, resource.WithProcessExecutableName())
		case DetectorProcessExecutablePath:
			opts = append(opts, resource.WithProcessExecutablePath())
		case DetectorProcessCommandArgs:
			opts = append(opts, resource.WithProcessCommandArgs())
		case DetectorProcessOwner:
			opts = append(opts, resource.WithProcessOwner())
		case DetectorProcessRuntimeName:
			opts = append(opts, resource.WithProcessRuntimeName())
		case DetectorProcessRuntimeVersion:
			opts = append(opts, resource.WithProcessRuntimeVersion())
		case DetectorProcessRuntimeDescription:
			opts = append(opts, resource.WithProcessRuntimeDescription())
		case DetectorContainer:
			opts = append(opts, resource.WithContainer())
		case DetectorContainerID:
			opts = append(opts, resource.WithContainerID())
		}
	}
	r, err := resource.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("resource: %w", err)
	}
	return r, nil
}

func newResourceAll(ctx context.Context) (*resource.Resource, error) {
	r, err := resource.New(ctx, []resource.Option{
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithHostID(),
		resource.WithTelemetrySDK(),
		resource.WithSchemaURL(os.Getenv(string(ENVResourceSchemaURL))),
		resource.WithOS(),
		resource.WithOSType(),
		resource.WithOSDescription(),
		resource.WithProcess(),
		resource.WithProcessPID(),
		resource.WithProcessExecutableName(),
		resource.WithProcessExecutablePath(),
		resource.WithProcessCommandArgs(),
		resource.WithProcessOwner(),
		resource.WithProcessRuntimeName(),
		resource.WithProcessRuntimeVersion(),
		resource.WithProcessRuntimeDescription(),
		resource.WithContainer(),
		resource.WithContainerID(),
	}...)
	if err != nil {
		return nil, fmt.Errorf("resource: %w", err)
	}
	return r, nil
}
