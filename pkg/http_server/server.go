package httpserver

import (
	"net"
	"net/http"
)

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
