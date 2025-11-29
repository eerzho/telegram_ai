package otelresource

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/resource"
)

func MustNew(ctx context.Context, cfg Config) *resource.Resource {
	res, err := New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return res
}

func New(ctx context.Context, cfg Config) (*resource.Resource, error) {
	var opts []resource.Option

	for _, detector := range cfg.Detectors {
		switch detector {
		case DetectorEnv:
			opts = append(opts, resource.WithFromEnv())
		case DetectorHost:
			opts = append(opts, resource.WithHost())
		case DetectorContainer:
			opts = append(opts, resource.WithContainer())
		}
	}

	r, err := resource.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("resource: %w", err)
	}

	return r, nil
}
