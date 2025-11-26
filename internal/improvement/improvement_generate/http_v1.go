package improvementgenerate

import (
	"log/slog"
	"net/http"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/eerzho/telegram-ai/pkg/json"
	"github.com/eerzho/telegram-ai/pkg/sse"
)

func HTTPv1(logger *slog.Logger, usecase *Usecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		ctx := r.Context()

		input, err := json.Decode[Input](r)
		if err != nil {
			logger.ErrorContext(ctx, "failed to json decode", slog.Any("error", err))
			json.EncodeError(w, err)
			return
		}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			logger.Log(ctx, domain.LogLevel(err), "failed to generate improvement", slog.Any("error", err))
			json.EncodeError(w, domain.MapToJSONError(err))
			return
		}

		sseWriter, err := sse.NewWriter(w)
		if err != nil {
			logger.ErrorContext(ctx, "failed to create sse writer")
			json.EncodeError(w, err)
			return
		}
		defer sseWriter.Close()

		if err = sseWriter.StreamFrom(ctx, &output); err != nil {
			logger.ErrorContext(ctx, "stream error", slog.Any("error", err))
		}
	})
}
