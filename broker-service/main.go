package main

import (
	"broker/internal"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	_ = godotenv.Load(".env")
	port := os.Getenv("BROKER_SERVICE_PORT")

	app := &internal.App{
		Cfg: &internal.Config{},
	}

	slog.Info("Starting broker service on port: %s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), app.Routes())
	if err != nil {
		log.Panic(err)
	}
}
