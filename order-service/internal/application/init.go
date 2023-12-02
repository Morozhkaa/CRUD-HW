// The application package is responsible for pre-creating the database and service and starting the application.
package application

import (
	"context"
	"fmt"
	"order-service/internal/adapters/db"
	"order-service/internal/adapters/http"
	"order-service/internal/domain/usecases"
	"time"
)

type App struct {
	opts          AppOptions
	shutdownFuncs []func(ctx context.Context) error
}

type AppOptions struct {
	DB_url      string
	HTTP_port   int
	Timeout     time.Duration
	IdleTimeout time.Duration
	AuthURL     string
	BillingURL  string
}

// New returns a new application instance.
func New(opts AppOptions) *App {
	return &App{
		opts: opts,
	}
}

// Start creates the database and service instances, then builds and starts the adapter.
func (app *App) Start() error {
	// creates the database and service instances
	storage, err := db.New(context.Background(), app.opts.DB_url)
	if err != nil {
		return fmt.Errorf("storage creation failed: %w", err)
	}
	userService := usecases.New(storage)

	// instantiate the adapter
	optsAdapter := http.AdapterOptions{
		HTTP_port:   app.opts.HTTP_port,
		Timeout:     app.opts.Timeout,
		IdleTimeout: app.opts.IdleTimeout,
		AuthURL:     app.opts.AuthURL,
		BillingURL:  app.opts.BillingURL,
	}
	s, err := http.New(userService, optsAdapter)
	if err != nil {
		return fmt.Errorf("adapter initialization failed: %w", err)
	}

	// add the application stop function to the list of shutdown functions and start the service.
	app.shutdownFuncs = append(app.shutdownFuncs, s.Stop)
	err = s.Start()
	if err != nil {
		return fmt.Errorf("server start failed: %w", err)
	}

	return nil
}

// Stop executes all shutdown functions.
func (a *App) Stop(ctx context.Context) error {
	var err error
	for i := len(a.shutdownFuncs) - 1; i >= 0; i-- {
		err = a.shutdownFuncs[i](ctx)
		if err != nil {
			return err
		}
	}
	return nil
}