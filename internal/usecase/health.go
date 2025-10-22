package usecase

import (
	"context"

	"github.com/eerzho/telegram-ai/internal/usecase/input"
	"github.com/eerzho/telegram-ai/internal/usecase/output"
)

type Health struct {
	version string
}

func NewHealth(
	version string,
) *Health {
	return &Health{
		version: version,
	}
}

func (h *Health) Check(ctx context.Context, input input.HealthCheck) (output.HealthCheck, error) {
	return output.HealthCheck{
		Status:  "ok",
		Version: h.version,
	}, nil
}
