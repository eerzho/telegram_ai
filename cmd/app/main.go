package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eerzho/telegram-ai/config"
	"github.com/eerzho/telegram-ai/internal/container"
	"github.com/eerzho/telegram-ai/internal/handler"
	"github.com/eerzho/telegram-ai/pkg/httpserver"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	c := container.New()

	cfg := c.MustGet("config").(config.Config)
	lgr := c.MustGet("logger").(*slog.Logger)

	mux := http.NewServeMux()
	handler.AddRoutes(mux, c)

	srv := httpserver.New(mux, cfg.HTTPServer)

	serverErrs := make(chan error, 1)
	go func() {
		lgr.Info("starting http server", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrs <- fmt.Errorf("ListenAndServe: %w", err)
		}
	}()

	select {
	case err := <-serverErrs:
		lgr.Error("http server error", slog.Any("error", err))
		return err
	case <-ctx.Done():
		lgr.Info("shutting down http server")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		lgr.Error("http server error", slog.Any("error", err))
		return err
	}

	c.MustReset()

	lgr.Info("http server stopped")

	return nil
}
