package providers

import (
	"context"
	"fmt"
	"github.com/SZabrodskii/microservices/libs/providers/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MongoDBProvider struct {
	client *mongo.Database
	logger *zap.Logger
}

func NewMongoDBProvider(cfg config.MongoDBConfig) (*MongoDBProvider, error) {

	clientOptions := options.Client().ApplyURI(cfg.GetDSN()).SetAuth(options.Credential{
		Username: cfg.GetUsername(),
		Password: cfg.GetPassword(),
	})

	logger := cfg.GetLogger()

	conn, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.Error("error connecting to MongoDB", zap.Error(err))
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	err = conn.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	provider := &MongoDBProvider{
		client: conn.Database(cfg.GetDatabase()),
		logger: logger,
	}

	return provider, nil
}

func (m *MongoDBProvider) Client() *mongo.Database {
	return m.client
}

func ProvideMongoDB() fx.Option {
	return fx.Provide(NewMongoDBProvider)
}
