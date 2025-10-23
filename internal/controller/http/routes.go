package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/eerzho/simpledi"
	"github.com/eerzho/telegram-ai/internal/usecase"
	"github.com/google/uuid"
)

func Handler(c *simpledi.Container) http.Handler {
	mux := http.NewServeMux()
	logger := c.MustGet("logger").(*slog.Logger)
	healthUsecase := c.MustGet("healthUsecase").(*usecase.Health)

	mux.Handle("GET /_hc", healthCheck(healthUsecase))
	mux.Handle("/", http.NotFoundHandler())

	var handler http.Handler = mux
	handler = loggingMiddleware(logger)(mux)
	handler = recoveryMiddleware(logger)(handler)
	return handler
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = uuid.NewString()
			}
			w.Header().Set("X-Request-Id", requestID)

			l := logger.With(
				slog.String("method", r.Method),
				slog.String("url_path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.String("request_id", requestID),
				slog.String("user_agent", r.UserAgent()),
				slog.Int("request_size", int(r.ContentLength)),
			)

			l.InfoContext(r.Context(), "request started")

			rw := &responseWriter{ResponseWriter: w}
			next.ServeHTTP(rw, r)

			l.InfoContext(r.Context(), "request finished",
				slog.Int("response_size", rw.size),
				slog.Int("status_code", rw.statusCode),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

func recoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.ErrorContext(r.Context(), "panic recovered",
						slog.Any("error", err),
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
					)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	size       int
	statusCode int
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}

func (r *responseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}
