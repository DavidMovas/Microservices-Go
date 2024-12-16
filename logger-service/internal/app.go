package internal

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"net/http"
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

func (a *App) Serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.Port),
		Handler: a.routes(),
	}

	return srv.ListenAndServe()
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)

	defer cancel()
	return a.db.Disconnect(ctx)
}

func (a *App) Insert(log Log) error {
	collection := a.db.Database("logs").Collection("logs")

	log.CreatedAt = time.Now()

	_, err := collection.InsertOne(context.Background(), log)
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

func (a *App) GetByID(ctx context.Context, id string) (*Log, error) {
	coll := a.db.Database("logs").Collection("logs")

	var log Log
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&log)
	return &log, err
}

func (a *App) GetLast(ctx context.Context, amount int) ([]*Log, error) {
	collection := a.db.Database("logs").Collection("logs")

	opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetLimit(int64(amount))
	cursor, err := collection.Find(ctx, bson.D{{"created_at", -1}}, opts)

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
