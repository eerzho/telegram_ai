package healthcheck

import (
	"context"

	"github.com/eerzho/telegram-ai/internal/config"
)

type Usecase struct {
	cfg config.App
}

func NewUsecase(
	cfg config.App,
) *Usecase {
	return &Usecase{
		cfg: cfg,
	}
}

func (h *Usecase) Execute(_ context.Context, _ Input) (Output, error) {
	return Output{
		Status:  "ok",
		Version: h.cfg.Version,
	}, nil
}
