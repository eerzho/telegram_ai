package http

import (
	"log/slog"
	"net/http"

	"github.com/eerzho/telegram-ai/internal/usecase"
	"github.com/eerzho/telegram-ai/internal/usecase/input"
	"github.com/eerzho/telegram-ai/pkg/json"
	"github.com/eerzho/telegram-ai/pkg/sse"
)

func streamAnswer(logger *slog.Logger, streamUsecase *usecase.Stream) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in, err := json.Decode[input.StreamAnswer](r)
		if err != nil {
			logger.ErrorContext(ctx, "failed to json decode", slog.Any("error", err))
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

		out, err := streamUsecase.Answer(ctx, in)
		if err != nil {
			logger.ErrorContext(ctx, "failed to answer", slog.Any("error", err))
			json.EncodeError(w, r, http.StatusInternalServerError, err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				logger.InfoContext(ctx, "client disconnected")
				return
			case err := <-out.ErrChan:
				if err != nil {
					logger.ErrorContext(ctx, "failed to answer", slog.Any("error", err))
					if err := sseWriter.Write(sse.Event{Name: "stop"}); err != nil {
						logger.WarnContext(ctx, "failed to write", slog.Any("error", err))
					}
					return
				}
			case text, ok := <-out.TextChan:
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
