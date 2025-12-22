package improvementgenerate

import (
	"context"
	"encoding/json"
	"time"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/eerzho/telegram_ai/internal/domain"
	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/semaphore"
)

const (
	generationTimeout = 20 // second
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
		return Output{}, errorhelp.WithOP(op, err)
	}

	if ok := u.sem.TryAcquire(1); !ok {
		return Output{}, errorhelp.WithOP(op, domain.ErrTooManyGenerateRequests)
	}

	textChan := make(chan string)
	errChan := make(chan error)
	go func() {
		defer u.sem.Release(1)
		defer close(textChan)
		defer close(errChan)

		genCtx, cancel := context.WithTimeoutCause(ctx, generationTimeout*time.Second, domain.ErrGenerationTimeout)
		defer cancel()

		err := u.generator.GenerateImprovement(genCtx, input.Text,
			func(chunk string) error {
				jsonChunk, err := json.Marshal(map[string]string{"text": chunk})
				if err != nil {
					return errorhelp.WithOP(op, err)
				}
				select {
				case <-genCtx.Done():
					return genCtx.Err()
				case textChan <- string(jsonChunk):
					return nil
				}
			},
		)
		if err != nil {
			select {
			case <-genCtx.Done():
			case errChan <- errorhelp.WithOP(op, err):
			}
			return
		}
	}()
	return Output{
		TextChan: textChan,
		ErrChan:  errChan,
	}, nil
}
