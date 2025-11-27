package responsegenerate

import (
	"fmt"
	"log/slog"
	"net/http"

	httphandler "github.com/eerzho/telegram-ai/pkg/http_handler"
	"github.com/eerzho/telegram-ai/pkg/json"
	"github.com/eerzho/telegram-ai/pkg/sse"
)

func HTTPv1(logger *slog.Logger, usecase *Usecase) httphandler.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "response_generate.HTTPv1"

		defer r.Body.Close()
		ctx := r.Context()

		input, err := json.Decode[Input](r)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		sseWriter, err := sse.NewWriter(w)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		defer sseWriter.Close()

		if err = sseWriter.StreamFrom(ctx, &output); err != nil {
			logger.ErrorContext(ctx, "stream error", slog.Any("error", fmt.Errorf("%s: %w", op, err)))
		}

		return nil
	}
}
