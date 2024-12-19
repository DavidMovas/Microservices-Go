package main

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
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

	esClient, err := elastic.NewClient(elastic.SetURL(cfg.EsURL), elastic.SetSniff(false))
	if err != nil {
		fmt.Printf("failed to create elastic client %v", err)
		os.Exit(1)
	}

	defer esClient.Stop()

	hook, err := elogrus.NewAsyncElasticHook(esClient, "localhost", logrus.InfoLevel, "logger")
	if err != nil {
		fmt.Printf("faildet to create elogrus hook: %v", err)
	}

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.AddHook(hook)

	app := internal.NewApp(cfg)

	log.Info("Starting logger service on port", "PORT", cfg.Port)

	if err = app.Serve(); err != nil {
		slog.Error("failed to start server", "ERROR", err)
		os.Exit(1)
	}
}
