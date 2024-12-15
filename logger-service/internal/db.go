package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	*mongo.Client
}

func NewMongoClient(config *Config) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(config.MongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: config.MongoUser,
		Password: config.MongoPassword,
	})

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return &MongoDB{Client: client}, nil
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
