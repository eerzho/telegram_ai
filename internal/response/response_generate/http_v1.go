package responsegenerate

import (
	"log/slog"
	"net/http"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	httpjson "github.com/eerzho/goiler/pkg/http_json"
	httpserver "github.com/eerzho/goiler/pkg/http_server"
	"github.com/eerzho/telegram_ai/pkg/sse"
)

// HTTPv1 godoc
//
// @summary response generate
// @tags response
//
// @accept json
// @param request body Input true "body"
//
// @produce json,event-stream
// @success 200 {object} sse.Event
// @failure 400 {object} httpjson.Error
// @failure 500 {object} httpjson.Error
//
// @router /v1/responses/generate [post].
func HTTPv1(logger *slog.Logger, usecase *Usecase) httpserver.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "response_generate.HTTPv1"

		ctx := r.Context()
		defer r.Body.Close()

		input, err := httpjson.Decode[Input](r)
		if err != nil {
			return errorhelp.WithOP(op, err)
		}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			return errorhelp.WithOP(op, err)
		}

		if err = sse.Stream(ctx, w, &output); err != nil {
			logger.ErrorContext(ctx, "stream error", slog.Any("error", errorhelp.WithOP(op, err)))
		}

		return nil
	}
}
