package http

import (
	"log/slog"
	"net/http"

	bodysize "github.com/eerzho/goiler/pkg/body_size"
	httpserver "github.com/eerzho/goiler/pkg/http_server"
	"github.com/eerzho/goiler/pkg/logging"
	"github.com/eerzho/goiler/pkg/recovery"
	"github.com/eerzho/simpledi"
	"github.com/eerzho/telegram_ai/internal/config"
	improvementgenerate "github.com/eerzho/telegram_ai/internal/improvement/improvement_generate"
	healthcheck "github.com/eerzho/telegram_ai/internal/monitoring/health_check"
	responsegenerate "github.com/eerzho/telegram_ai/internal/response/response_generate"
	summarygenerate "github.com/eerzho/telegram_ai/internal/summary/summary_generate"
	"github.com/eerzho/telegram_ai/pkg/cors"
	swagger "github.com/swaggo/http-swagger"
)

// Handler godoc
//
// @version 1.0
// @title TelegramAI
// @description Telegram with AI features
//
// @schemes http
// @host localhost
// @basePath /
//
// @externalDocs.description GitHub
// @externalDocs.url https://github.com/eerzho/telegram_ai
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", swagger.WrapHandler)

	cfg := simpledi.Get[config.Config]("config")
	logger := simpledi.Get[*slog.Logger]("logger")
	errorHandler := errorHandler(logger)

	mux.Handle(
		"GET /_hc",
		httpserver.Wrap(healthcheck.HTTPv1(
			simpledi.Get[*healthcheck.Usecase]("healthCheckUsecase"),
		), errorHandler),
	)

	mux.Handle(
		"POST /v1/responses/generate",
		httpserver.Wrap(responsegenerate.HTTPv1(
			logger,
			simpledi.Get[*responsegenerate.Usecase]("responseGenerateUsecase"),
		), errorHandler),
	)

	mux.Handle(
		"POST /v1/summaries/generate",
		httpserver.Wrap(summarygenerate.HTTPv1(
			logger,
			simpledi.Get[*summarygenerate.Usecase]("summaryGenerateUsecase"),
		), errorHandler),
	)

	mux.Handle(
		"POST /v1/improvements/generate",
		httpserver.Wrap(improvementgenerate.HTTPv1(
			logger,
			simpledi.Get[*improvementgenerate.Usecase]("improvementGenerateUsecase"),
		), errorHandler),
	)

	var handler http.Handler = mux
	handler = cors.Middleware(cfg.CORS)(handler)
	handler = bodysize.Middleware(cfg.BodySize)(handler)
	handler = logging.Middleware(logger)(handler)
	handler = recovery.Middleware(logger)(handler)
	return handler
}
