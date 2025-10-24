package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/eerzho/telegram-ai/internal/usecase/input"
	"github.com/eerzho/telegram-ai/internal/usecase/output"
	"github.com/go-playground/validator/v10"
)

type Stream struct {
	logger   *slog.Logger
	validate *validator.Validate
}

func NewStream(
	logger *slog.Logger,
	validate *validator.Validate,
) *Stream {
	return &Stream{
		logger:   logger,
		validate: validate,
	}
}

func (s *Stream) Answer(ctx context.Context, in input.StreamAnswer) (output.StreamAnswer, error) {
	const op = "usecase.Stream.Answer"

	if err := s.validate.Struct(in); err != nil {
		return output.StreamAnswer{}, fmt.Errorf("%s: %w", op, err)
	}

	textChan := make(chan string, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(textChan)
		defer close(errChan)

		for _, msg := range in.Messages {
			textChan <- msg.Text
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return output.StreamAnswer{
		TextChan: textChan,
		ErrChan:  errChan,
	}, nil
}
