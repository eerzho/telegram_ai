package cors

import (
	"net/http"
	"strconv"
	"strings"
)

type Config struct {
	AllowedOrigins   string `env:"CORS_ALLOWED_ORIGINS"   envDefault:"*"`
	AllowedMethods   string `env:"CORS_ALLOWED_METHODS"   envDefault:"GET,POST,PUT,DELETE,OPTIONS,PATCH"`
	AllowedHeaders   string `env:"CORS_ALLOWED_HEADERS"   envDefault:"Content-Type,Authorization,X-Request-Id"`
	AllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS" envDefault:"false"`
	MaxAge           int    `env:"CORS_MAX_AGE"           envDefault:"3600"`
}

func Middleware(cfg Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && isOriginAllowed(origin, cfg.AllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if cfg.AllowedOrigins == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			if cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", cfg.AllowedMethods)
				w.Header().Set("Access-Control-Allow-Headers", cfg.AllowedHeaders)
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isOriginAllowed(origin, allowedOrigins string) bool {
	if allowedOrigins == "*" {
		return true
	}
	origins := strings.SplitSeq(allowedOrigins, ",")
	for allowed := range origins {
		allowed = strings.TrimSpace(allowed)
		if allowed == origin {
			return true
		}
	}
	return false
}
