package main

import (
	"auth/internal"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	_ = godotenv.Load(".env")
	port := os.Getenv("AUTH_SERVICE_PORT")

	db, err := connectDB(os.Getenv("POSTGRES_URL"))
	if err != nil {
		slog.Error("failed to connect to database: %w\n", err)
		os.Exit(1)
	}

	app := &internal.App{
		DB: db,
	}

	slog.Info("Starting auth service on port: %s\n", port)

	if err = http.ListenAndServe(fmt.Sprintf(":%s", port), app.Routes()); err != nil {
		log.Panic(err)
	}
}

func connectDB(connSting string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := pgxpool.New(ctx, connSting)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
