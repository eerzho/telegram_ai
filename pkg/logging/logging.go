package logging

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/eerzho/telegram-ai/pkg/logger"
	"github.com/google/uuid"
)

func Middleware(lgr *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = uuid.NewString()
				r.Header.Set("X-Request-Id", requestID)
			}
			w.Header().Set("X-Request-Id", requestID)
			rw := &responseWriter{ResponseWriter: w}

			ctx := context.WithValue(r.Context(), logger.RequestIDKey, requestID)
			r = r.WithContext(ctx)

			next.ServeHTTP(rw, r)

			lgr.InfoContext(r.Context(), "http request",
				slog.String("method", r.Method),
				slog.String("host", r.Host),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("url.scheme", r.URL.Scheme),
				slog.String("url.path", r.URL.Path),
				slog.String("url.query", r.URL.RawQuery),
				slog.String("request.header.user_agent", r.UserAgent()),
				slog.String("request.header.request_id", requestID),
				slog.String("request.header.content_type", r.Header.Get("Content-Type")),
				slog.Int64("request.size", r.ContentLength),
				slog.Int("response.size", rw.size),
				slog.Int("response.status_code", rw.statusCode),
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
