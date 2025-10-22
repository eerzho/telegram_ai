package http

import (
	"net/http"

	"github.com/eerzho/telegram-ai/internal/usecase"
	"github.com/eerzho/telegram-ai/internal/usecase/input"
	"github.com/eerzho/telegram-ai/pkg/json"
)

func healthCheck(healthUsecase *usecase.Health) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		output, err := healthUsecase.Check(ctx, input.HealthCheck{})
		if err != nil {
			json.EncodeError(w, r, http.StatusInternalServerError, err)
			return
		}

		json.Encode(w, r, http.StatusOK, output)
	})
}
