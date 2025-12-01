package otelruntime

import "time"

type Config struct {
	MinimumReadMemStatsInterval time.Duration `env:"OTEL_METRIC_RUNTIME_MIN_READ_MEMSTATS_INTERVAL" envDefault:"15s"`
}
