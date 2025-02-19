package main

import (
	"callbox/internal/app"
	"log/slog"
	"os"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := app.New(logger)

	if err := app.Start(); err != nil {
		logger.Error("Failed to start server")
	}
}
