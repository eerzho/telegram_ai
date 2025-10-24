package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/eerzho/telegram-ai/internal/usecase/input"
	"github.com/eerzho/telegram-ai/internal/usecase/output"
)

type Stream struct {
	logger *slog.Logger
}

func NewStream(logger *slog.Logger) *Stream {
	return &Stream{
		logger: logger,
	}
}

func (s *Stream) Answer(ctx context.Context, in input.StreamAnswer) output.StreamAnswer {
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
	}
}
