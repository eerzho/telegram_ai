package summary_get

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/go-playground/validator/v10"
)

type Valkey interface {
	SetSummary(ctx context.Context, chatID, text string) error
	GetSummary(ctx context.Context, chatID string) (string, error)
}

type Postgres interface {
	FindSummary(ctx context.Context, chatID string) (domain.Summary, error)
}

type Usecase struct {
	logger   *slog.Logger
	validate *validator.Validate
	valkey   Valkey
	postgres Postgres
}

func NewUsecase(
	logger *slog.Logger,
	validate *validator.Validate,
	valkey Valkey,
	postgres Postgres,
) *Usecase {
	return &Usecase{
		logger:   logger,
		validate: validate,
		valkey:   valkey,
		postgres: postgres,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "summary_get.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, fmt.Errorf("%s: %w", op, err)
	}

	text, err := u.valkey.GetSummary(ctx, input.ChatID)
	if err == nil {
		return Output{Text: text}, nil
	}
	u.logger.WarnContext(ctx, "failed to get summary",
		slog.Any("error", fmt.Errorf("%s: %w", op, err)),
	)

	summary, err := u.postgres.FindSummary(ctx, input.ChatID)
	if err != nil {
		return Output{}, fmt.Errorf("%s: %w", op, err)
	}

	err = u.valkey.SetSummary(ctx, summary.ChatID, summary.Text)
	if err != nil {
		u.logger.WarnContext(ctx, "failed to set summary",
			slog.Any("error", fmt.Errorf("%s: %w", op, err)),
		)
	}

	return Output{Text: text}, nil
}
