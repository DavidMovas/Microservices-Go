package internal

import (
	"fmt"
	"net/http"
	"time"
)

const (
	GracefulShutdownTimeout = 5 * time.Second
)

type App struct {
	config *Config
}

func NewApp(config *Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.Port),
		Handler: a.routes(),
	}

	return srv.ListenAndServe()
}
