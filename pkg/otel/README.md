# OpenTelemetry Configuration

--- 

## Resource Configuration

**Configuration** (`otel_resource/config.go`):

```
OTEL_RESOURCE_DETECTORS=env,host,container
```

Supported values: `env`, `host`, `container`

For resource configuration, see the official documentation:
[OpenTelemetry](https://opentelemetry.io/docs/concepts/sdk-configuration/general-sdk-configuration/)

---

## Trace Exporter Configuration

**Configuration** (`otel_trace_exporter/config.go`):

```
OTEL_TRACES_EXPORTER=otlp-grpc
```

Supported values: `otlp-grpc`, `otlp-http`, `console`, `none`

For trace exporter configuration, see the official documentation:
[OpenTelemetry](https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/)

---

## Tracer Configuration

**Configuration** (`otel_tracer/config.go`):

```
OTEL_BSP_SCHEDULE_DELAY=5s
OTEL_BSP_EXPORT_TIMEOUT=30s
OTEL_BSP_MAX_EXPORT_BATCH_SIZE=512
OTEL_BSP_MAX_QUEUE_SIZE=2048
```

For racer configuration, see the official documentation:
[OpenTelemetry](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#batch-span-processor)

---

## Metric Exporter Configuration

**Configuration** (`otel_metric_exporter/config.go`):

```
OTEL_METRICS_EXPORTER=otlp-grpc
```

Supported values: `otlp-grpc`, `otlp-http`, `console`, `none`

For metric exporter configuration, see the official documentation:
[OpenTelemetry](https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/)

---

## Meter Configuration

**Configuration** (`otel_meter/config.go`):

```
OTEL_METRIC_CARDINALITY_LIMIT=2000
```

For meter configuration, see the official documentation:
[OpenTelemetry](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#periodic-metric-reader)

---

## Metric Runtime Configuration

**Configuration** (`otel_metric_runtime/config.go`):

```
OTEL_METRIC_RUNTIME_MIN_READ_MEMSTATS_INTERVAL=15s
```

For runtime metrics configuration, see the official documentation:
[OpenTelemetry](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/runtime)

---

## Log Exporter Configuration

**Configuration** (`otel_log_exporter/config.go`):

```
OTEL_LOGS_EXPORTER=otlp-grpc
```

Supported values: `otlp-grpc`, `otlp-http`, `console`, `none`

For log exporter configuration, see the official documentation:
[OpenTelemetry](https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/)

---

## Logger Configuration

**Configuration** (`otel_logger/logger.go`):

For logger configuration, see the official documentation:
[OpenTelemetry](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#batch-logrecord-processor)
