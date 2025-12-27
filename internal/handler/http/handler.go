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
	generateimprovement "github.com/eerzho/telegram_ai/internal/improvement/generate_improvement"
	healthcheck "github.com/eerzho/telegram_ai/internal/monitoring/health_check"
	generateresponse "github.com/eerzho/telegram_ai/internal/response/generate_response"
	createsetting "github.com/eerzho/telegram_ai/internal/setting/create_setting"
	deletesetting "github.com/eerzho/telegram_ai/internal/setting/delete_setting"
	getsetting "github.com/eerzho/telegram_ai/internal/setting/get_setting"
	updatesetting "github.com/eerzho/telegram_ai/internal/setting/update_setting"
	generatesummary "github.com/eerzho/telegram_ai/internal/summary/generate_summary"
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
		httpserver.Wrap(generateresponse.HTTPv1(
			logger,
			simpledi.Get[*generateresponse.Usecase]("responseGenerateUsecase"),
		), errorHandler),
	)

	mux.Handle(
		"POST /v1/summaries/generate",
		httpserver.Wrap(generatesummary.HTTPv1(
			logger,
			simpledi.Get[*generatesummary.Usecase]("summaryGenerateUsecase"),
		), errorHandler),
	)

	mux.Handle(
		"POST /v1/improvements/generate",
		httpserver.Wrap(generateimprovement.HTTPv1(
			logger,
			simpledi.Get[*generateimprovement.Usecase]("improvementGenerateUsecase"),
		), errorHandler),
	)

	mux.Handle(
		"POST /v1/settings",
		httpserver.Wrap(createsetting.HTTPv1(
			simpledi.Get[*createsetting.Usecase]("createSettingUsecase"),
		), errorHandler),
	)

	mux.Handle(
		"DELETE /v1/settings/{user_id}/{chat_id}",
		httpserver.Wrap(deletesetting.HTTPv1(
			simpledi.Get[*deletesetting.Usecase]("deleteSettingUsecase"),
		), errorHandler),
	)

	mux.Handle(
		"PUT /v1/settings/{user_id}/{chat_id}",
		httpserver.Wrap(updatesetting.HTTPv1(
			simpledi.Get[*updatesetting.Usecase]("updateSettingUsecase"),
		), errorHandler),
	)

	mux.Handle(
		"GET /v1/settings/{user_id}/{chat_id}",
		httpserver.Wrap(getsetting.HTTPv1(
			simpledi.Get[*getsetting.Usecase]("getSettingUsecase"),
		), errorHandler),
	)

	var handler http.Handler = mux
	handler = cors.Middleware(cfg.CORS)(handler)
	handler = bodysize.Middleware(cfg.BodySize)(handler)
	handler = logging.Middleware(logger)(handler)
	handler = recovery.Middleware(logger)(handler)
	return handler
}
