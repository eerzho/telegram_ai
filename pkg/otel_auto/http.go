package otelauto

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewHandler(handler http.Handler) http.Handler {
	opts := make([]otelhttp.Option, 0)

	opts = append(opts, otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
		return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
	}))

	return otelhttp.NewHandler(handler, "", opts...)
}
