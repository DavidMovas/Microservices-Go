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
	_ = app

	slog.Info("Starting logger service on port", "PORT", cfg.Port)

}
