package generatesummary

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
	generationTimeout = 30 // second
)

type generator interface {
	GenerateSummary(
		ctx context.Context,
		language string,
		dialog domain.Dialog,
		onChunk func(chunk string) error,
	) error
}

type Usecase struct {
	sem       *semaphore.Weighted
	validate  *validator.Validate
	generator generator
}

func NewUsecase(
	sem *semaphore.Weighted,
	validate *validator.Validate,
	generator generator,
) *Usecase {
	return &Usecase{
		sem:       sem,
		validate:  validate,
		generator: generator,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "summary_generate.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, errorhelp.WithOP(op, err)
	}

	ok := u.sem.TryAcquire(1)
	if !ok {
		return Output{}, errorhelp.WithOP(op, domain.ErrTooManyGenerateRequests)
	}

	textChan := make(chan string)
	errChan := make(chan error)
	go func() {
		defer u.sem.Release(1)
		defer close(textChan)
		defer close(errChan)

		genCtx, genCancel := context.WithTimeoutCause(ctx, generationTimeout*time.Second, domain.ErrGenerationTimeout)
		defer genCancel()

		err := u.generator.GenerateSummary(genCtx, input.Language, input.ToDialog(),
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
