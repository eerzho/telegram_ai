package response_generate

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/semaphore"
)

type Generator interface {
	GenerateResponse(
		ctx context.Context,
		dialog domain.Dialog,
		onChunk func(chunk string) error,
	) error
}

type Usecase struct {
	sem       *semaphore.Weighted
	validate  *validator.Validate
	generator Generator
}

func NewUsecase(
	sem *semaphore.Weighted,
	validate *validator.Validate,
	generator Generator,
) *Usecase {
	return &Usecase{
		sem:       sem,
		validate:  validate,
		generator: generator,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "response_generate.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, fmt.Errorf("%s: %w", op, err)
	}

	ok := u.sem.TryAcquire(1)
	if !ok {
		return Output{}, fmt.Errorf("%s: %w", op, domain.ErrTooManyGenerateRequests)
	}

	textChan := make(chan string, 10)
	errChan := make(chan error, 1)
	go func() {
		defer u.sem.Release(1)
		defer close(textChan)
		defer close(errChan)

		slices.SortFunc(input.Messages, func(a, b InputMessage) int {
			return cmp.Compare(a.Date, b.Date)
		})
		dialog := inputToDialog(input)

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
			select {
			case errChan <- fmt.Errorf("%s: %w", op, err):
			case <-ctx.Done():
			}
			return
		}
	}()

	return Output{
		TextChan: textChan,
		ErrChan:  errChan,
	}, nil
}
