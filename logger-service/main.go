package main

import (
	"log/slog"
	"logger/internal"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := internal.NewConfig()
	if err != nil {
		slog.Error("failed to parse config", "ERROR", err)
		os.Exit(1)
	}

	mdb, err := internal.NewMongoClient(cfg)
	if err != nil {
		slog.Error("failed to connect to mongo", "ERROR", err)
		os.Exit(1)
	}

	app := internal.NewApp(cfg, mdb)

	slog.Info("Starting logger service on port", "PORT", cfg.Port)

	defer func() {
		err = app.Shutdown()
		if err != nil {
			slog.Error("failed to shutdown", "ERROR", err)
			os.Exit(1)
		}
	}()

	if err = app.Serve(); err != nil {
		slog.Error("failed to start server", "ERROR", err)
		os.Exit(1)
	}
}
