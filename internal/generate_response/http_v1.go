package generate_response

import (
	"log/slog"
	"net/http"

	"github.com/eerzho/telegram-ai/pkg/json"
	"github.com/eerzho/telegram-ai/pkg/sse"
)

func HTTPv1(logger *slog.Logger, usecase *Usecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		input, err := json.Decode[Input](r)
		if err != nil {
			logger.ErrorContext(ctx, "failed to json decode", slog.Any("error", err))
			json.EncodeError(w, r, http.StatusBadRequest, err)
			return
		}

		output, err := usecase.Execute(ctx, input)
		if err != nil {
			logger.ErrorContext(ctx, "failed to generate response", slog.Any("error", err))
			json.EncodeError(w, r, http.StatusInternalServerError, err)
			return
		}

		sseWriter, err := sse.NewWriter(w)
		if err != nil {
			logger.ErrorContext(ctx, "failed to create sse writer", slog.Any("error", err))
			json.EncodeError(w, r, http.StatusInternalServerError, err)
			return
		}
		defer func() {
			err := sseWriter.Close()
			if err != nil {
				logger.WarnContext(ctx, "failed to close sse writer", slog.Any("error", err))
			}
		}()

		if err := sseWriter.Write(sse.Event{Name: "start"}); err != nil {
			logger.WarnContext(ctx, "failed to write", slog.Any("error", err))
			return
		}

		for {
			select {
			case <-ctx.Done():
				logger.InfoContext(ctx, "client disconnected")
				return
			case err := <-output.ErrChan:
				if err != nil {
					logger.ErrorContext(ctx, "failed to generate response", slog.Any("error", err))
					if err := sseWriter.Write(sse.Event{Name: "stop"}); err != nil {
						logger.WarnContext(ctx, "failed to write", slog.Any("error", err))
					}
					return
				}
			case text, ok := <-output.TextChan:
				if !ok {
					if err := sseWriter.Write(sse.Event{Name: "stop"}); err != nil {
						logger.WarnContext(ctx, "failed to write", slog.Any("error", err))
					}
					return
				}
				if err := sseWriter.Write(sse.Event{Name: "append", Data: text}); err != nil {
					logger.WarnContext(ctx, "failed to write", slog.Any("error", err))
					return
				}
			}
		}
	})
}
