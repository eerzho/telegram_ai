package http

import (
	"log/slog"
	"net/http"

	"github.com/eerzho/simpledi"
	"github.com/eerzho/telegram-ai/config"
	healthcheck "github.com/eerzho/telegram-ai/internal/health/health_check"
	improvementgenerate "github.com/eerzho/telegram-ai/internal/improvement/improvement_generate"
	responsegenerate "github.com/eerzho/telegram-ai/internal/response/response_generate"
	summarygenerate "github.com/eerzho/telegram-ai/internal/summary/summary_generate"
	bodysize "github.com/eerzho/telegram-ai/pkg/body_size"
	"github.com/eerzho/telegram-ai/pkg/cors"
	httphandler "github.com/eerzho/telegram-ai/pkg/http_handler"
	"github.com/eerzho/telegram-ai/pkg/logging"
	"github.com/eerzho/telegram-ai/pkg/recovery"
	swagger "github.com/swaggo/http-swagger"
)

// Handler godoc
//
// @title TelegramAI API
//
// @schemes http
// @host localhost
// @basePath /
//
// @externalDocs.description GitHub
// @externalDocs.url https://github.com/eerzho/telegram-ai
func Handler() http.Handler {
	mux := http.NewServeMux()
	cfg := simpledi.Get[config.Config]("config")
	logger := simpledi.Get[*slog.Logger]("logger")

	errorHandler := errorHandler(logger)

	mux.Handle("/swagger/", swagger.WrapHandler)

	mux.Handle(
		"GET /_hc",
		httphandler.Wrap(healthcheck.HTTPv1(
			simpledi.Get[*healthcheck.Usecase]("healthCheckUsecase"),
		), errorHandler),
	)
	mux.Handle(
		"POST /v1/responses/generate",
		httphandler.Wrap(responsegenerate.HTTPv1(
			logger,
			simpledi.Get[*responsegenerate.Usecase]("responseGenerateUsecase"),
		), errorHandler),
	)
	mux.Handle(
		"POST /v1/summaries/generate",
		httphandler.Wrap(summarygenerate.HTTPv1(
			logger,
			simpledi.Get[*summarygenerate.Usecase]("summaryGenerateUsecase"),
		), errorHandler),
	)
	mux.Handle(
		"POST /v1/improvements/generate",
		httphandler.Wrap(improvementgenerate.HTTPv1(
			logger,
			simpledi.Get[*improvementgenerate.Usecase]("improvementGenerateUsecase"),
		), errorHandler),
	)

	var handler http.Handler = mux
	handler = bodysize.Middleware(cfg.BodySize)(handler)
	handler = cors.Middleware(cfg.CORS)(handler)
	handler = logging.Middleware(logger)(handler)
	handler = recovery.Middleware(logger)(handler)
	return handler
}
