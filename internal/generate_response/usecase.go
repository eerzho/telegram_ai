package generate_response

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Genkit interface {
	GenerateResponse(
		ctx context.Context,
		dialog string,
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

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("My name is %s", input.Owner.Name))
	for _, msg := range input.Messages {
		if input.Owner.ChatID == msg.Sender.ChatID {
			sb.WriteString(fmt.Sprintf("\nI said: %s", msg.Text))
		} else {
			sb.WriteString(fmt.Sprintf("\n%s said: %s", msg.Sender.Name, msg.Text))
		}
	}

	textChan := make(chan string, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(textChan)
		defer close(errChan)

		err := s.genkit.GenerateResponse(ctx, sb.String(), func(chunk string) error {
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
