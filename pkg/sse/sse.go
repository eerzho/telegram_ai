package sse

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
)

const (
	defaultBufferSize = 4096
)

var (
	ErrStreamingNotSupported = errors.New("streaming not supported")
)

type Writer struct {
	w       *bufio.Writer
	flusher http.Flusher
}

func NewWriter(w http.ResponseWriter) (*Writer, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("sse: %w", ErrStreamingNotSupported)
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	return &Writer{
		w:       bufio.NewWriterSize(w, defaultBufferSize),
		flusher: flusher,
	}, nil
}

func (s *Writer) Write(event, data string) error {
	_, err := fmt.Fprintf(s.w, "event: %s\ndata: %s\n\n", event, data)
	if err != nil {
		return fmt.Errorf("sse write: %w", err)
	}
	return s.flush()
}

func (s *Writer) Close() error {
	return s.flush()
}

func (s *Writer) flush() error {
	if err := s.w.Flush(); err != nil {
		return fmt.Errorf("sse flush: %w", err)
	}
	s.flusher.Flush()
	return nil
}
