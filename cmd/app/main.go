package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/eerzho/setting/config"
	"github.com/eerzho/setting/pkg/logger"
)

func main() {
	c, err := config.Init()
	if err != nil {
		panic(err)
	}

	l, err := logger.Init(c.Logger)
	if err != nil {
		panic(err)
	}

	//
	// run application
}

func AppRun(ctx context.Context, c config.Config) error {
	// init simpledi
	//
	// setup router
	//
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	// destroy simpledi
	//
	return nil
}
