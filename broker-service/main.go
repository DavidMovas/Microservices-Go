package main

import (
	"broker/internal"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := internal.NewConfig()
	if err != nil {
		slog.Error("failed to parse config", "ERROR", err)
		os.Exit(1)
	}

	app := &internal.App{
		Cfg: cfg,
	}

	slog.Info("Starting broker service on port", "PORT", cfg.Port)

	if err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), app.Routes()); err != nil {
		log.Panic(err)
	}
}
