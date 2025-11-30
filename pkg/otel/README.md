# OpenTelemetry Configuration

--- 

## Resource Configuration

**Custom configuration** (`otel_resource/config.go`):

```
OTEL_RESOURCE_DETECTORS=env,host,container
```

Supported values: `env`, `host`, `container`

For all other resource configuration, see the official documentation:
📖 [OpenTelemetry Resource Configuration](https://opentelemetry.io/docs/concepts/sdk-configuration/general-sdk-configuration/)

---

## Trace Exporter Configuration

**Custom configuration** (`otel_trace_exporter/config.go`):

```
OTEL_TRACES_EXPORTER=otlp-grpc
```

Supported values: `otlp-grpc`, `otlp-http`, `console`, `none`

For all other trace exporter configuration, see the official documentation:
📖 [OpenTelemetry OTLP Exporter Configuration](https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/)

---

## Tracer Configuration

**Custom configuration** (`otel_tracer/config.go`):

```
OTEL_BSP_SCHEDULE_DELAY=5s
OTEL_BSP_EXPORT_TIMEOUT=30s
OTEL_BSP_MAX_EXPORT_BATCH_SIZE=512
OTEL_BSP_MAX_QUEUE_SIZE=2048
```

For all other tracer configuration, see the official documentation:
📖 [OpenTelemetry Batch Span Processor Configuration](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#batch-span-processor)

---

## Metric Exporter Configuration

**Custom configuration** (`otel_metric_exporter/config.go`):

```
OTEL_METRICS_EXPORTER=otlp-grpc
```

Supported values: `otlp-grpc`, `otlp-http`, `console`, `none`

For all other metric exporter configuration, see the official documentation:
📖 [OpenTelemetry OTLP Exporter Configuration](https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/)

---

## Meter Configuration

**Custom configuration** (`otel_meter/config.go`):

```
OTEL_METRIC_CARDINALITY_LIMIT=2000
```

For all other meter configuration, see the official documentation:
📖 [OpenTelemetry Periodic Metric Reader Configuration](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#periodic-metric-reader)

---

## Metric Runtime Configuration

**Custom configuration** (`otel_metric_runtime/config.go`):

```
OTEL_METRIC_RUNTIME_MIN_READ_MEMSTATS_INTERVAL=15s
```

For all other runtime metrics configuration, see the official documentation:
📖 [OpenTelemetry Runtime Instrumentation](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/runtime)
