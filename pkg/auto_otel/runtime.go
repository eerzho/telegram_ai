package autootel

import (
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

func runtimeStart() error {
	opts := make([]runtime.Option, 0)
	if envInterval := os.Getenv("OTEL_RUNTIME_METRICS_INTERVAL"); envInterval != "" {
		if duration, err := time.ParseDuration(envInterval); err == nil {
			opts = append(opts, runtime.WithMinimumReadMemStatsInterval(duration))
		}
	}
	if err := runtime.Start(opts...); err != nil {
		return fmt.Errorf("runtime: %w", err)
	}
	return nil
}
