package generate

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/go-playground/validator/v10"
)

type Generator interface {
	GenerateSummary(
		ctx context.Context,
		language string,
		dialog domain.Dialog,
		onChunk func(chunk string) error,
	) error
}

type Cacher interface {
	SetSummary(ctx context.Context, chatID, summary string) error
}

type Usecase struct {
	validate  *validator.Validate
	generator Generator
	cacher    Cacher
}

func NewUsecase(
	validate *validator.Validate,
	generator Generator,
	cacher Cacher,
) *Usecase {
	return &Usecase{
		validate:  validate,
		generator: generator,
		cacher:    cacher,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "generate_summary.Usecase.Execute"

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

		summary := ""
		err := u.generator.GenerateSummary(ctx, input.Language, dialog,
			func(chunk string) error {
				summary += chunk
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

		err = u.cacher.SetSummary(ctx, input.Owner.ChatID, summary)
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
