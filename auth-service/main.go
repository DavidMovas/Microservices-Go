package main

import (
	"auth/internal"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const dbTimeout = time.Second * 10

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := internal.NewConfig()
	if err != nil {
		slog.Error("failed to parse config", "ERROR", err)
		os.Exit(1)
	}

	db, err := connectDB(cfg.ConnString)
	if err != nil {
		slog.Error("failed to connect to database", "ERROR", err)
		os.Exit(1)
	}

	app := &internal.App{
		DB:  db,
		Cfg: cfg,
	}

	slog.Info("Starting auth service on port", "PORT", cfg.Port)

	if err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), app.Routes()); err != nil {
		log.Panic(err)
	}
}

func connectDB(connSting string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	slog.Info("Connecting to database", "CONNECTION_STRING", connSting)
	db, err := pgxpool.New(ctx, connSting)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	slog.Info("Connected to database")
	return db, nil
}
