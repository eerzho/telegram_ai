package summary_generate

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
			json.EncodeError(w, r, err)
			return
		}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			logger.Log(ctx, domain.LogLevel(err), "failed to generate summary", slog.Any("error", err))
			json.EncodeError(w, r, domain.MapToJSONError(err))
			return
		}

		sseWriter, err := sse.NewWriter(w)
		if err != nil {
			logger.ErrorContext(ctx, "failed to create sse writer", slog.Any("error", err))
			json.EncodeError(w, r, err)
			return
		}
		defer sseWriter.Close()

		if err := sseWriter.Write(ctx, sse.Event{Name: "start"}); err != nil {
			logger.ErrorContext(ctx, "failed to write", slog.Any("error", err))
			return
		}

		for {
			select {
			case text, ok := <-output.TextChan:
				if !ok {
					if err := sseWriter.Write(ctx, sse.Event{Name: "stop"}); err != nil {
						logger.ErrorContext(ctx, "failed to write", slog.Any("error", err))
					}
					return
				}
				if err := sseWriter.Write(ctx, sse.Event{Name: "append", Data: text}); err != nil {
					logger.ErrorContext(ctx, "failed to write", slog.Any("error", err))
					_ = sseWriter.Write(ctx, sse.Event{Name: "stop_with_error", Data: "Please try again later."})
					return
				}
			case err := <-output.ErrChan:
				if err != nil {
					logger.ErrorContext(ctx, "failed to generate summary", slog.Any("error", err))
					_ = sseWriter.Write(ctx, sse.Event{Name: "stop_with_error", Data: "Please try again later."})
					return
				}
			}
		}
	})
}
