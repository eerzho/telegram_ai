package handler

import (
	"net/http"

	"github.com/eerzho/simpledi"
	v1 "github.com/eerzho/telegram-ai/internal/handler/v1"
)

func AddRoutes(mux *http.ServeMux, c *simpledi.Container) {
	health := c.MustGet("healthHandler").(*v1.Health)
	stream := c.MustGet("streamHandler").(*v1.Stream)

	mux.HandleFunc("/healths/check", health.Check)
	mux.HandleFunc("/streams/answer", stream.Answer)
}
