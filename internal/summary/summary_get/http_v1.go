package summary_get

import (
	"log/slog"
	"net/http"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/eerzho/telegram-ai/pkg/json"
)

func HTTPv1(logger *slog.Logger, usecase *Usecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		input := Input{
			OwnerID: r.PathValue("owner_id"),
			PeerID:  r.PathValue("peer_id"),
		}
		output, err := usecase.Execute(ctx, input)
		if err != nil {
			logger.Log(ctx, domain.LogLevel(err),
				"failed to get summary",
				slog.Any("error", err),
			)
			json.EncodeError(w, r, domain.MapToJSONError(err))
			return
		}

		json.Encode(w, r, http.StatusOK, output)
	})
}
