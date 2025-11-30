package oteltracing

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "http-server"

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tracer := otel.Tracer(tracerName)

			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = uuid.NewString()
				r.Header.Set("X-Request-Id", requestID)
			}
			w.Header().Set("X-Request-Id", requestID)

			ctx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.String()),
				trace.WithAttributes(
					attribute.String("http.request.method", r.Method),
					attribute.String("url.scheme", r.URL.Scheme),
					attribute.String("url.path", r.URL.Path),
					attribute.String("url.query", r.URL.RawQuery),
					attribute.String("server.address", r.Host),
					attribute.String("client.address", r.RemoteAddr),
					attribute.String("user_agent.original", r.UserAgent()),
					attribute.String("http.request.header.request_id", requestID),
					attribute.String("http.request.header.content_type", r.Header.Get("Content-Type")),
					attribute.Int64("http.request.size", r.ContentLength),
				),
			)
			defer span.End()

			rw := &responseWriter{ResponseWriter: w}

			r = r.WithContext(ctx)
			next.ServeHTTP(rw, r)

			span.SetAttributes(
				attribute.Int("http.response.status_code", rw.statusCode),
				attribute.Int("http.response.size", rw.size),
			)

			if rw.statusCode >= 400 {
				span.SetStatus(codes.Error, "HTTP error")
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
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
