package generate

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/go-playground/validator/v10"
)

type Generator interface {
	GenerateResponse(
		ctx context.Context,
		dialog domain.Dialog,
		onChunk func(chunk string) error,
	) error
}

type Usecase struct {
	validate  *validator.Validate
	generator Generator
}

func NewUsecase(
	validate *validator.Validate,
	generator Generator,
) *Usecase {
	return &Usecase{
		validate:  validate,
		generator: generator,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "generate_response.Usecase.Execute"

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

		err := u.generator.GenerateResponse(ctx, dialog,
			func(chunk string) error {
				select {
				case textChan <- chunk:
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
	}()

	return Output{
		TextChan: textChan,
		ErrChan:  errChan,
	}, nil
}
