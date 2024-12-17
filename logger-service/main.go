package main

import (
	"github.com/sirupsen/logrus"
	"log/slog"
	"logger/internal"
	"os"
)

func main() {
	cfg, err := internal.NewConfig()
	if err != nil {
		slog.Error("failed to parse config", "ERROR", err)
		os.Exit(1)
	}

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.AddHook(internal.NewLokiHook(cfg.LokiConfig))

	app := internal.NewApp(cfg)

	log.Info("Starting logger service on port", "PORT", cfg.Port)

	if err = app.Serve(); err != nil {
		slog.Error("failed to start server", "ERROR", err)
		os.Exit(1)
	}
}
