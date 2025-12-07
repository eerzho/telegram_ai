package autootel

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
	opts := buildResourceOptions(detectors)
	r, err := resource.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("resource: %w", err)
	}
	return r, nil
}

func buildResourceOptions(detectors []string) []resource.Option {
	opts := make([]resource.Option, 0, len(detectors))
	for _, detector := range detectors {
		opt := getDetectorOption(strings.TrimSpace(detector))
		if opt != nil {
			opts = append(opts, opt)
		}
	}
	return opts
}

func getDetectorOption(detector string) resource.Option {
	switch detector {
	case DetectorENV:
		return resource.WithFromEnv()
	case DetectorHost:
		return resource.WithHost()
	case DetectorHostID:
		return resource.WithHostID()
	case DetectorTelemetrySDK:
		return resource.WithTelemetrySDK()
	case DetectorSchemaURL:
		return resource.WithSchemaURL(os.Getenv(ENVResourceSchemaURL))
	case DetectorOS:
		return resource.WithOS()
	case DetectorOSType:
		return resource.WithOSType()
	case DetectorOSDescription:
		return resource.WithOSDescription()
	case DetectorProcess:
		return resource.WithProcess()
	case DetectorProcessPID:
		return resource.WithProcessPID()
	case DetectorProcessExecutableName:
		return resource.WithProcessExecutableName()
	case DetectorProcessExecutablePath:
		return resource.WithProcessExecutablePath()
	case DetectorProcessCommandArgs:
		return resource.WithProcessCommandArgs()
	case DetectorProcessOwner:
		return resource.WithProcessOwner()
	case DetectorProcessRuntimeName:
		return resource.WithProcessRuntimeName()
	case DetectorProcessRuntimeVersion:
		return resource.WithProcessRuntimeVersion()
	case DetectorProcessRuntimeDescription:
		return resource.WithProcessRuntimeDescription()
	case DetectorContainer:
		return resource.WithContainer()
	case DetectorContainerID:
		return resource.WithContainerID()
	default:
		return nil
	}
}

func newResourceAll(ctx context.Context) (*resource.Resource, error) {
	r, err := resource.New(ctx, []resource.Option{
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithHostID(),
		resource.WithTelemetrySDK(),
		resource.WithSchemaURL(os.Getenv(ENVResourceSchemaURL)),
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
