package bodysize

import (
	"net/http"
)

const (
	MB = 1 << 20 // 1048576 bytes
)

type Config struct {
	Max int `env:"BODY_SIZE_MAX" envDefault:"5"`
}

func Middleware(cfg Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			maxBytes := int64(cfg.Max) * MB
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
