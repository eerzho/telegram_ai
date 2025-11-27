package healthcheck

import (
	"log/slog"
	"net/http"

	errorhelp "github.com/eerzho/telegram-ai/pkg/error_help"
	httphandler "github.com/eerzho/telegram-ai/pkg/http_handler"
	"github.com/eerzho/telegram-ai/pkg/json"
)

func HTTPv1(logger *slog.Logger, usecase *Usecase) httphandler.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "health_check.HTTPv1"

		defer r.Body.Close()
		ctx := r.Context()

		output, err := usecase.Execute(ctx, Input{})
		if err != nil {
			return errorhelp.WithOP(op, err)
		}

		json.Encode(w, http.StatusOK, output)
		return nil
	}
}
