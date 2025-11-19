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

func Handler() http.Handler {
	mux := http.NewServeMux()
	cfg := simpledi.Get[config.Config]("config")
	logger := simpledi.Get[*slog.Logger]("logger")
	healthCheckUsecase := simpledi.Get[*health_check.Usecase]("healthCheckUsecase")
	generateResponseUsecase := simpledi.Get[*generate_response.Usecase]("generateResponseUsecase")
	generateSummaryUsecase := simpledi.Get[*generate_summary.Usecase]("generateSummaryUsecase")

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
