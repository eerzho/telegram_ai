package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/eerzho/goiler/docs"
	autootel "github.com/eerzho/goiler/pkg/auto_otel"
	httpserver "github.com/eerzho/goiler/pkg/http_server"
	"github.com/eerzho/simpledi"
	"github.com/eerzho/telegram_ai/internal/config"
	"github.com/eerzho/telegram_ai/internal/container"
	"github.com/eerzho/telegram_ai/internal/handler/http"
)

const (
	shutdownTimeout = 10 // second
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	signalCtx, signalCancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signalCancel()

	otel, err := autootel.Setup(signalCtx)
	if err != nil {
		return fmt.Errorf("failed to setup otel: %w", err)
	}

	for _, definition := range container.Definitions() {
		simpledi.Set(definition)
	}

	simpledi.Resolve()

	cfg := simpledi.Get[config.Config]("config")
	lgr := simpledi.Get[*slog.Logger]("logger")

	httpServer := httpserver.New(
		autootel.NewHandler(http.Handler()),
		cfg.HTTPServer,
	)

	serverErrs := make(chan error, 1)
	go func() {
		lgr.Info("starting http server", slog.String("addr", httpServer.Addr))
		if err = httpServer.ListenAndServe(); err != nil {
			serverErrs <- err
		}
	}()

	select {
	case err = <-serverErrs:
		return fmt.Errorf("server: %w", err)
	case <-signalCtx.Done():
		lgr.InfoContext(signalCtx, "shutdown signal received")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, shutdownTimeout*time.Second)
	defer shutdownCancel()

	if err = httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	if err = simpledi.Close(); err != nil {
		lgr.WarnContext(shutdownCtx, "failed to close container", slog.Any("error", err))
	}

	lgr.InfoContext(shutdownCtx, "http server stopped")

	if err = otel.Shutdown(shutdownCtx); err != nil {
		lgr.WarnContext(shutdownCtx, "failed to shutdown otel", slog.Any("error", err))
	}

	return nil
}
