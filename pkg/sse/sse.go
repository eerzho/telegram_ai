package sse

import (
	"bufio"
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
)

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

func (s *Writer) Write(e Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e.ID != 0 {
		if _, err := fmt.Fprintf(s.w, "id: %d\n", e.ID); err != nil {
			return fmt.Errorf("sse write: %w", err)
		}
	}
	if e.Name != "" {
		if _, err := fmt.Fprintf(s.w, "event: %s\n", e.Name); err != nil {
			return fmt.Errorf("sse write: %w", err)
		}
	}
	if e.Data != "" {
		if _, err := fmt.Fprintf(s.w, "data: %s\n", e.Data); err != nil {
			return fmt.Errorf("sse write: %w", err)
		}
	}
	if e.Retry != 0 {
		if _, err := fmt.Fprintf(s.w, "retry: %d\n", e.Retry); err != nil {
			return fmt.Errorf("sse write: %w", err)
		}
	}
	if _, err := s.w.WriteString("\n"); err != nil {
		return fmt.Errorf("sse write: %w", err)
	}
	return s.flush()
}

func (s *Writer) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flush()
}

func (s *Writer) flush() error {
	if err := s.w.Flush(); err != nil {
		return fmt.Errorf("sse flush: %w", err)
	}
	s.f.Flush()
	return nil
}
