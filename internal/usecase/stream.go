package usecase

import (
	"context"

	"github.com/eerzho/telegram-ai/internal/usecase/input"
	"github.com/eerzho/telegram-ai/internal/usecase/output"
)

type Stream struct {
}

func NewStream() *Stream {
	return &Stream{}
}

func (s *Stream) Answer(ctx context.Context, intput input.StreamAnswer) (output.StreamAnswer, <-chan error) {
	resultChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errChan)
	}()
}
