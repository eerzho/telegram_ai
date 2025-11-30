package otelmetricruntime

import (
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

func MustNew(cfg Config) {
	if err := New(cfg); err != nil {
		panic(err)
	}
}

func New(cfg Config) error {
	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(cfg.MinimumReadMemStatsInterval)); err != nil {
		return fmt.Errorf("runtime metrics: %w", err)
	}
	return nil
}
