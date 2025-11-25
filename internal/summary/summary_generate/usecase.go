package summary_generate

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/semaphore"
)

type Generator interface {
	GenerateSummary(
		ctx context.Context,
		language string,
		dialog domain.Dialog,
		onChunk func(chunk string) error,
	) error
}

type Usecase struct {
	sem       *semaphore.Weighted
	logger    *slog.Logger
	validate  *validator.Validate
	generator Generator
}

func NewUsecase(
	sem *semaphore.Weighted,
	logger *slog.Logger,
	validate *validator.Validate,
	generator Generator,
) *Usecase {
	return &Usecase{
		sem:       sem,
		logger:    logger,
		validate:  validate,
		generator: generator,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "summary_generate.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, fmt.Errorf("%s: %w", op, err)
	}

	ok := u.sem.TryAcquire(1)
	if !ok {
		return Output{}, fmt.Errorf("%s: %w", op, domain.ErrTooManyGenerateRequests)
	}

	slices.SortFunc(input.Messages, func(a, b InputMessage) int {
		return cmp.Compare(a.Date, b.Date)
	})
	dialog := inputToDialog(input)

	textChan := make(chan string)
	errChan := make(chan error)
	go func() {
		defer u.sem.Release(1)
		defer close(textChan)
		defer close(errChan)

		genCtx, genCancel := context.WithTimeoutCause(ctx, 30*time.Second, domain.ErrGenerationTimeout)
		defer genCancel()

		var builder strings.Builder
		err := u.generator.GenerateSummary(genCtx, input.Language, dialog,
			func(chunk string) error {
				builder.WriteString(chunk)
				jsonChunk, err := json.Marshal(map[string]string{"text": chunk})
				if err != nil {
					return fmt.Errorf("%s: %w", op, err)
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
