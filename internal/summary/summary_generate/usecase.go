package summary_generate

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/go-playground/validator/v10"
)

type Genkit interface {
	GenerateSummary(
		ctx context.Context,
		language string,
		dialog domain.Dialog,
		onChunk func(chunk string) error,
	) error
}

type Valkey interface {
	SetSummary(ctx context.Context, chatID, text string) error
}

type Postgres interface {
	UpdateSummary(ctx context.Context, chatID, text string) error
}

type Usecase struct {
	logger   *slog.Logger
	validate *validator.Validate
	genkit   Genkit
	postgres Postgres
	valkey   Valkey
}

func NewUsecase(
	logger *slog.Logger,
	validate *validator.Validate,
	genkit Genkit,
	postgres Postgres,
	valkey Valkey,
) *Usecase {
	return &Usecase{
		logger:   logger,
		validate: validate,
		genkit:   genkit,
		postgres: postgres,
		valkey:   valkey,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "summary_generate.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, fmt.Errorf("%s: %w", op, err)
	}

	slices.SortFunc(input.Messages, func(a, b InputMessage) int {
		return cmp.Compare(a.Date, b.Date)
	})
	dialog := inputToDialog(input)

	textChan := make(chan string, 10)
	errChan := make(chan error, 1)
	go func() {
		defer close(textChan)
		defer close(errChan)

		text := ""
		err := u.genkit.GenerateSummary(ctx, input.Language, dialog,
			func(chunk string) error {
				text += chunk
				data, err := json.Marshal(map[string]string{"text": chunk})
				if err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
				select {
				case textChan <- string(data):
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			},
		)
		if err != nil {
			errChan <- fmt.Errorf("%s: %w", op, err)
			return
		}

		err = u.postgres.UpdateSummary(ctx, input.Owner.ChatID, text)
		if err != nil {
			errChan <- fmt.Errorf("%s: %w", op, err)
			return
		}

		err = u.valkey.SetSummary(ctx, input.Owner.ChatID, text)
		if err != nil {
			u.logger.WarnContext(ctx, "failed to set summary",
				slog.Any("error", fmt.Errorf("%s: %w", op, err)),
			)
		}
	}()

	return Output{
		TextChan: textChan,
		ErrChan:  errChan,
	}, nil
}
