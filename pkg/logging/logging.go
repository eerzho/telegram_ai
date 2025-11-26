package logging

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func Middleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = uuid.NewString()
			}

			w.Header().Set("X-Request-Id", requestID)
			rw := &responseWriter{ResponseWriter: w}

			ctx := context.WithValue(r.Context(), "request_id", requestID)
			r = r.WithContext(ctx)

			next.ServeHTTP(rw, r)

			logger.InfoContext(r.Context(), "http request",
				slog.String("method", r.Method),
				slog.String("url_path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.String("user_agent", r.UserAgent()),
				slog.Int("request_size", int(r.ContentLength)),
				slog.Int("response_size", rw.size),
				slog.Int("status_code", rw.statusCode),
				slog.Duration("duration", time.Since(start)),
			)
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

func (r *responseWriter) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
