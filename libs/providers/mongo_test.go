package providers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.uber.org/zap"
)

func TestMongoDBConnection(t *testing.T) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://172.17.0.3:27017"))
	assert.NoError(t, err, "failed to connect to MongoDB")
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	assert.NoError(t, err, "failed to ping MongoDB")

	assert.NotNil(t, client, "MongoDB client is nil")
}
