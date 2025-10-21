package v1

import "net/http"

type Stream struct {
}

func NewStream() *Stream {
	return &Stream{}
}

func (s *Stream) Answer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// todo: call openapi
}
