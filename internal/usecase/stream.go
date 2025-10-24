package usecase

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
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

	slices.SortFunc(in.Messages, func(a, b input.StreamAnswerMessage) int {
		return cmp.Compare(a.Date, b.Date)
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("My name is %s", in.Owner.Name))
	for _, msg := range in.Messages {
		if in.Owner.ChatID == msg.Sender.ChatID {
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
