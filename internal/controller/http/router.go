package http

import (
	"log/slog"
	"net/http"

	"github.com/eerzho/simpledi"
	"github.com/eerzho/telegram-ai/config"
	health_check "github.com/eerzho/telegram-ai/internal/health/check"
	response_generate "github.com/eerzho/telegram-ai/internal/response/generate"
	summary_generate "github.com/eerzho/telegram-ai/internal/summary/generate"
	"github.com/eerzho/telegram-ai/pkg/cors"
	"github.com/eerzho/telegram-ai/pkg/logging"
	"github.com/eerzho/telegram-ai/pkg/recovery"
)

func Handler() http.Handler {
	mux := http.NewServeMux()
	cfg := simpledi.Get[config.Config]("config")
	logger := simpledi.Get[*slog.Logger]("logger")

	healthCheckUsecase := simpledi.Get[*health_check.Usecase]("healthCheckUsecase")
	responseGenerateUsecase := simpledi.Get[*response_generate.Usecase]("responseGenerateUsecase")
	summaryGenerateUsecase := simpledi.Get[*summary_generate.Usecase]("summaryGenerateUsecase")

	mux.Handle("GET /v1/health-check", health_check.HTTPv1(logger, healthCheckUsecase))
	mux.Handle("POST /v1/responses/generate", response_generate.HTTPv1(logger, responseGenerateUsecase))
	mux.Handle("POST /v1/summaries/generate", summary_generate.HTTPv1(logger, summaryGenerateUsecase))
	mux.Handle("/", http.NotFoundHandler())

	var handler http.Handler = mux
	handler = cors.Middleware(cfg.CORS)(handler)
	handler = logging.Middleware(logger)(handler)
	handler = recovery.Middleware(logger)(handler)
	return handler
}
