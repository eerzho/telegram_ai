package httpserver

import (
	"net"
	"net/http"
	"time"
)

type Config struct {
	Host              string        `env:"HTTP_SERVER_HOST"          envDefault:""`
	Port              string        `env:"HTTP_SERVER_PORT"          envDefault:"8080"`
	ReadHeaderTimeout time.Duration `env:"HTTP_SERVER_READ_HEADER_TIMEOUT" envDefault:"10s"`
	ReadTimeout       time.Duration `env:"HTTP_SERVER_READ_TIMEOUT"  envDefault:"10s"`
	WriteTimeout      time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout       time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT"  envDefault:"60s"`
}

func New(handler http.Handler, cfg Config) *http.Server {
	return &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
}
