package sse

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

const (
	defaultBufferSize = 4096
)

var (
	ErrStreamingNotSupported = errors.New("streaming not supported")
	ErrClientDisconnected    = errors.New("client disconnected")
)

type Streamer interface {
	Next() (Event, bool)
}

type Writer struct {
	w  *bufio.Writer
	f  http.Flusher
	mu sync.Mutex
}

type Event struct {
	ID    int
	Name  string
	Data  string
	Retry int // milliseconds
}

func Stream(ctx context.Context, w http.ResponseWriter, s Streamer) error {
	writer, err := newWriter(w)
	if err != nil {
		return err
	}
	defer writer.close()
	return writer.stream(ctx, s)
}

func newWriter(w http.ResponseWriter) (*Writer, error) {
	const op = "sse.newWriter"
	f, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, ErrStreamingNotSupported)
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	return &Writer{
		w: bufio.NewWriterSize(w, defaultBufferSize),
		f: f,
	}, nil
}

func (w *Writer) stream(ctx context.Context, s Streamer) error {
	const op = "sse.Writer.stream"
	for {
		event, done := s.Next()
		if err := w.write(ctx, event); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if done {
			return nil
		}
	}
}

func (w *Writer) close() error {
	const op = "sse.Writer.close"
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.flush(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (w *Writer) write(ctx context.Context, e Event) error {
	const op = "sse.Writer.write"
	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ErrClientDisconnected)
	default:
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if e.ID != 0 {
		if _, err := fmt.Fprintf(w.w, "id: %d\n", e.ID); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	if e.Name != "" {
		if _, err := fmt.Fprintf(w.w, "event: %s\n", e.Name); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	if e.Data != "" {
		if _, err := fmt.Fprintf(w.w, "data: %s\n", e.Data); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	if e.Retry != 0 {
		if _, err := fmt.Fprintf(w.w, "retry: %d\n", e.Retry); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	if _, err := w.w.WriteString("\n"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return w.flush()
}

func (w *Writer) flush() error {
	const op = "sse.Writer.flush"
	if err := w.w.Flush(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	w.f.Flush()
	return nil
}
