package summary_get

import (
	"log/slog"
	"net/http"

	"github.com/eerzho/telegram-ai/pkg/json"
)

func HTTPv1(logger *slog.Logger, usecase *Usecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		input := Input{ChatID: r.PathValue("id")}
		output, err := usecase.Execute(ctx, input)
		if err != nil {
			logger.ErrorContext(ctx, "failed to get summary", slog.Any("error", err))
			json.EncodeError(w, r, http.StatusInternalServerError, err)
			return
		}

		json.Encode(w, r, http.StatusOK, output)
	})
}
