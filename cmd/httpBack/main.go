package main

import (
	"log/slog"
	"os"
	"track-selection/internal/bootstrap"
	"track-selection/internal/config"
)

func main() {
	cfg := config.MustLoadConfigApp()

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cfg.GetLogLevel(),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	app := bootstrap.NewApp(cfg)
	app.Run()
}
