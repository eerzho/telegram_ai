package generate_response

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/go-playground/validator/v10"
)

type Genkit interface {
	GenerateResponse(
		ctx context.Context,
		dialog domain.Dialog,
		onChunk func(chunk string) error,
	) error
}

type Usecase struct {
	logger   *slog.Logger
	validate *validator.Validate
	genkit   Genkit
}

func NewUsecase(
	logger *slog.Logger,
	validate *validator.Validate,
	genkit Genkit,
) *Usecase {
	return &Usecase{
		logger:   logger,
		validate: validate,
		genkit:   genkit,
	}
}

func (s *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "generate_response.Usecase.Execute"

	if err := s.validate.Struct(input); err != nil {
		return Output{}, fmt.Errorf("%s: %w", op, err)
	}

	slices.SortFunc(input.Messages, func(a, b InputMessage) int {
		return cmp.Compare(a.Date, b.Date)
	})

	textChan := make(chan string, 10)
	errChan := make(chan error, 1)
	go func() {
		defer close(textChan)
		defer close(errChan)

		dialog := inputToDialog(input)
		err := s.genkit.GenerateResponse(ctx, dialog, func(chunk string) error {
			select {
			case textChan <- chunk:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
		if err != nil {
			errChan <- fmt.Errorf("%s: %w", op, err)
			return
		}
	}()

	return Output{
		TextChan: textChan,
		ErrChan:  errChan,
	}, nil
}
