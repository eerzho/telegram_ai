package http

import (
	"log/slog"
	"net/http"

	"github.com/eerzho/simpledi"
	"github.com/eerzho/telegram-ai/config"
	"github.com/eerzho/telegram-ai/internal/usecase"
	"github.com/eerzho/telegram-ai/pkg/cors"
	"github.com/eerzho/telegram-ai/pkg/logging"
	"github.com/eerzho/telegram-ai/pkg/recovery"
)

func Handler(c *simpledi.Container) http.Handler {
	mux := http.NewServeMux()
	cfg := c.MustGet("config").(config.Config)
	logger := c.MustGet("logger").(*slog.Logger)
	healthUsecase := c.MustGet("healthUsecase").(*usecase.Health)
	streamUsecase := c.MustGet("streamUsecase").(*usecase.Stream)

	mux.Handle("GET /_hc", healthCheck(logger, healthUsecase))
	mux.Handle("POST /stream/answer", streamAnswer(logger, streamUsecase))
	mux.Handle("/", http.NotFoundHandler())

	var handler http.Handler = mux
	handler = cors.Middleware(cfg.CORS)(handler)
	handler = logging.Middleware(logger)(handler)
	handler = recovery.Middleware(logger)(handler)
	return handler
}
