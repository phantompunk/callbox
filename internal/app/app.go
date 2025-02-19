package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type app struct {
	Logger *slog.Logger
}

func New(logger *slog.Logger) *app {
	return &app{
		Logger: logger,
	}
}

func (app *app) Start() error {
	app.Logger.Info("Starting server at 4000")

	server := &http.Server{
		Addr:         ":4000",
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	done := make(chan bool, 1)

	go gracefulShutdown(app, server, done)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		app.Logger.Error("Failed to listen and serve", slog.Any("error", err))
	}

	<-done
	app.Logger.Info("Graceful shutdown complete")
	return nil
}

func gracefulShutdown(app *app, server *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	app.Logger.Info("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error %v", err)
	}

	app.Logger.Info("Server exiting")
	done <- true
}
