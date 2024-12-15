package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

const (
	GracefulShutdownTimeout = 5 * time.Second
)

type App struct {
	db     *MongoDB
	config *Config
}

func NewApp(config *Config, db *MongoDB) *App {
	return &App{
		config: config,
		db:     db,
	}
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)

	defer cancel()
	return a.db.Disconnect(ctx)
}

func (a *App) Insert(log Log) error {
	collection := a.db.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), log)
	return err
}

func (a *App) GetLogs(ctx context.Context) ([]*Log, error) {
	collection := a.db.Database("logs").Collection("logs")

	cursor, err := collection.Find(ctx, bson.D{{"created_at", -1}})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = cursor.Close(ctx)
	}()

	var logs []*Log
	err = cursor.All(ctx, &logs)
	return logs, err
}
