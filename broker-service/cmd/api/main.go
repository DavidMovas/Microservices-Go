package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "8010"

type App struct {
	cfg *Config
}

type Config struct{}

func main() {
	app := &App{
		cfg: &Config{},
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting broker service on port: %s\n", webPort)

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
