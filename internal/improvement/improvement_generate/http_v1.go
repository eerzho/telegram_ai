package improvementgenerate

import (
	"log/slog"
	"net/http"

	errorhelp "github.com/eerzho/telegram-ai/pkg/error_help"
	httphandler "github.com/eerzho/telegram-ai/pkg/http_handler"
	"github.com/eerzho/telegram-ai/pkg/json"
	"github.com/eerzho/telegram-ai/pkg/sse"
)

// HTTPv1 godoc
//
// @summary improvement generate
// @tags improvement
//
// @accept json
// @param request body Input true "body"
//
// @produce json,event-stream
// @success 200 {object} sse.Event
// @failure 400 {object} json.Error
// @failure 500 {object} json.Error
//
// @router /v1/improvements/generate [post].
func HTTPv1(logger *slog.Logger, usecase *Usecase) httphandler.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "improvement_generate.HTTPv1"

		defer r.Body.Close()
		ctx := r.Context()

		input, err := json.Decode[Input](r)
		if err != nil {
			return errorhelp.WithOP(op, err)
		}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			return errorhelp.WithOP(op, err)
		}

		sseWriter, err := sse.NewWriter(w)
		if err != nil {
			return errorhelp.WithOP(op, err)
		}
		defer sseWriter.Close()

		if err = sseWriter.Stream(ctx, &output); err != nil {
			logger.ErrorContext(ctx, "stream error", slog.Any("error", errorhelp.WithOP(op, err)))
		}

		return nil
	}
}
