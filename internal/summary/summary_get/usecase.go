package summary_get

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/go-playground/validator/v10"
)

type Cache interface {
	SetSummary(ctx context.Context, ownerID, peerID, text string) error
	GetSummary(ctx context.Context, ownerID, peerID string) (string, error)
}

type Storage interface {
	FindSummary(ctx context.Context, ownerID, peerID string) (domain.Summary, error)
}

type Usecase struct {
	logger   *slog.Logger
	validate *validator.Validate
	cache    Cache
	storage  Storage
}

func NewUsecase(
	logger *slog.Logger,
	validate *validator.Validate,
	cache Cache,
	storage Storage,
) *Usecase {
	return &Usecase{
		logger:   logger,
		validate: validate,
		cache:    cache,
		storage:  storage,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "summary_get.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, fmt.Errorf("%s: %w", op, err)
	}

	text, err := u.cache.GetSummary(ctx, input.OwnerID, input.PeerID)
	if err == nil {
		return Output{Text: text}, nil
	}
	u.logger.InfoContext(ctx, "failed to get summary",
		slog.Any("error", fmt.Errorf("%s: %w", op, err)),
	)

	summary, err := u.storage.FindSummary(ctx, input.OwnerID, input.PeerID)
	if err != nil {
		return Output{}, fmt.Errorf("%s: %w", op, err)
	}

	err = u.cache.SetSummary(ctx, summary.OwnerID, summary.PeerID, summary.Text)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to set summary",
			slog.Any("error", fmt.Errorf("%s: %w", op, err)),
		)
	}

	return Output{Text: summary.Text}, nil
}
