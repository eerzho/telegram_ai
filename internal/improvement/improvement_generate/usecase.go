package improvement_generate

import (
	"context"
	"fmt"
	"time"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/semaphore"
)

type Generator interface {
	GenerateImprovement(
		ctx context.Context,
		text string,
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
	const op = "improvement_generate.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, fmt.Errorf("%s: %w", op, err)
	}

	if ok := u.sem.TryAcquire(1); !ok {
		return Output{}, fmt.Errorf("%s: %w", op, domain.ErrTooManyGenerateRequests)
	}

	textChan := make(chan string)
	errChan := make(chan error)
	go func() {
		defer u.sem.Release(1)
		defer close(textChan)
		defer close(errChan)

		genCtx, cancel := context.WithTimeoutCause(ctx, 20*time.Second, domain.ErrGenerationTimeout)
		defer cancel()

		err := u.generator.GenerateImprovement(genCtx, input.Text,
			func(chunk string) error {
				select {
				case <-genCtx.Done():
					return genCtx.Err()
				case textChan <- chunk:
					return nil
				}
			},
		)
		if err != nil {
			select {
			case <-genCtx.Done():
			case errChan <- fmt.Errorf("%s: %w", op, err):
			}
			return
		}
	}()
	return Output{
		TextChan: textChan,
		ErrChan:  errChan,
	}, nil
}
