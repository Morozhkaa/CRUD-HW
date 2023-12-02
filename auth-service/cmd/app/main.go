// The main package is the entry point. Initializes the configuration file and logger, and runs
// the application's main logic. Performs graceful shutdown.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "user-service/api/swagger/public"
	"user-service/internal/application"
	"user-service/internal/config"
	"user-service/pkg/infra/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	// generate config from env vars
	cfg := config.Get()

	// initialize logger
	optsLogger := logger.LoggerOptions{IsProd: cfg.IsProd}
	log := logger.New(optsLogger)

	optsApp := application.AppOptions{
		DB_url:      cfg.DB_URL,
		HTTP_port:   cfg.HTTP_port,
		Timeout:     cfg.Timeout,
		IdleTimeout: cfg.IdleTimeout,
	}
	app := application.New(optsApp)

	err := app.Start()
	if err != nil {
		log.Error("app not started", "desc", err.Error())
	}
	<-ctx.Done()

	// graceful shutdown
	log.Info("shutting down...")
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()

	err = app.Stop(stopCtx)
	if err != nil {
		log.Error("app stop error", "desc", err.Error())
	}
	log.Info("app stopped")
}
