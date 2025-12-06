package otelauto

import (
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

func runtimeStart() error {
	interval := 15 * time.Second
	if envInterval := os.Getenv("OTEL_RUNTIME_METRICS_INTERVAL"); envInterval != "" {
		if duration, err := time.ParseDuration(envInterval); err == nil {
			interval = duration
		}
	}
	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(interval)); err != nil {
		return fmt.Errorf("runtime: %w", err)
	}
	return nil
}
