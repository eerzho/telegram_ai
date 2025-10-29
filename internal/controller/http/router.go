package http

import (
	"log/slog"
	"net/http"

	"github.com/eerzho/simpledi"
	"github.com/eerzho/telegram-ai/config"
	"github.com/eerzho/telegram-ai/internal/infra/health_check"
	"github.com/eerzho/telegram-ai/internal/response/generate_response"
	"github.com/eerzho/telegram-ai/internal/summary/generate_summary"
	"github.com/eerzho/telegram-ai/pkg/cors"
	"github.com/eerzho/telegram-ai/pkg/logging"
	"github.com/eerzho/telegram-ai/pkg/recovery"
)

func Handler(c *simpledi.Container) http.Handler {
	mux := http.NewServeMux()
	cfg := c.MustGet("config").(config.Config)
	logger := c.MustGet("logger").(*slog.Logger)
	healthCheckUsecase := c.MustGet("healthCheckUsecase").(*health_check.Usecase)
	generateResponseUsecase := c.MustGet("generateResponseUsecase").(*generate_response.Usecase)
	generateSummaryUsecase := c.MustGet("generateSummaryUsecase").(*generate_summary.Usecase)

	mux.Handle("GET /v1/health-check", health_check.HTTPv1(logger, healthCheckUsecase))
	mux.Handle("POST /v1/responses/generate", generate_response.HTTPv1(logger, generateResponseUsecase))
	mux.Handle("POST /v1/summaries/generate", generate_summary.HTTPv1(logger, generateSummaryUsecase))
	mux.Handle("/", http.NotFoundHandler())

	var handler http.Handler = mux
	handler = cors.Middleware(cfg.CORS)(handler)
	handler = logging.Middleware(logger)(handler)
	handler = recovery.Middleware(logger)(handler)
	return handler
}
