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

func NewWriter(w http.ResponseWriter) (*Writer, error) {
	f, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("sse: %w", ErrStreamingNotSupported)
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

func (w *Writer) Write(ctx context.Context, e Event) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("sse: %w", ErrClientDisconnected)
	default:
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if e.ID != 0 {
		if _, err := fmt.Fprintf(w.w, "id: %d\n", e.ID); err != nil {
			return fmt.Errorf("sse write: %w", err)
		}
	}
	if e.Name != "" {
		if _, err := fmt.Fprintf(w.w, "event: %s\n", e.Name); err != nil {
			return fmt.Errorf("sse write: %w", err)
		}
	}
	if e.Data != "" {
		if _, err := fmt.Fprintf(w.w, "data: %s\n", e.Data); err != nil {
			return fmt.Errorf("sse write: %w", err)
		}
	}
	if e.Retry != 0 {
		if _, err := fmt.Fprintf(w.w, "retry: %d\n", e.Retry); err != nil {
			return fmt.Errorf("sse write: %w", err)
		}
	}
	if _, err := w.w.WriteString("\n"); err != nil {
		return fmt.Errorf("sse write: %w", err)
	}
	return w.flush()
}

func (w *Writer) StreamFrom(ctx context.Context, s Streamer) error {
	for {
		event, done := s.Next()
		if err := w.Write(ctx, event); err != nil {
			return err
		}
		if done {
			return nil
		}
	}
}

func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.flush()
}

func (w *Writer) flush() error {
	if err := w.w.Flush(); err != nil {
		return fmt.Errorf("sse flush: %w", err)
	}
	w.f.Flush()
	return nil
}
