package logging

import (
	"log/slog"
	"net/http"
	"time"
)

func Middleware(lgr *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(rw, r)

			lgr.InfoContext(r.Context(), "request handled",
				slog.String("method", r.Method),
				slog.String("host", r.Host),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("scheme", r.URL.Scheme),
				slog.String("path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.String("user_agent", r.UserAgent()),
				slog.String("content_type", r.Header.Get("Content-Type")),
				slog.Int64("request_size", r.ContentLength),
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
