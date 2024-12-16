package main

import (
	"log/slog"
	"mail/internal"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := internal.NewConfig()
	if err != nil {
		slog.Error("failed to parse config", "ERROR", err)
		os.Exit(1)
	}

	app := internal.NewApp(cfg)

	defer func() {
		err = app.Shutdown()
		if err != nil {
			slog.Error("failed to shutdown", "ERROR", err)
			os.Exit(1)
		}
	}()

	slog.Info("Starting mail service on port", "PORT", cfg.Port)

	if err = app.Serve(); err != nil {
		slog.Error("failed to start server", "ERROR", err)
		os.Exit(1)
	}
}
