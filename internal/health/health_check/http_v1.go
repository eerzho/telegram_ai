package healthcheck

import (
	"log/slog"
	"net/http"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/eerzho/telegram-ai/pkg/json"
)

func HTTPv1(logger *slog.Logger, usecase *Usecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		output, err := usecase.Execute(ctx, Input{})
		if err != nil {
			logger.Log(ctx, domain.LogLevel(err),
				"failed to health check",
				slog.Any("error", err),
			)
			json.EncodeError(w, domain.MapToJSONError(err))
			return
		}

		json.Encode(w, http.StatusOK, output)
	})
}
