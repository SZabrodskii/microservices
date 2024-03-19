package config

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDBParams struct {
	DSN      string
	Username string
	Password string
	Database string
	Options  *options.ClientOptions
	Logger   *zap.Logger
}

type MongoDBConfig interface {
	GetDSN() string
	GetUsername() string
	GetPassword() string
	GetDatabase() string
	GetOptions() *options.ClientOptions
	GetLogger() *zap.Logger
}

func (mp *MongoDBParams) GetDSN() string {
	return mp.DSN
}

func (mp *MongoDBParams) GetUsername() string {
	return mp.Username
}

func (mp *MongoDBParams) GetPassword() string {
	return mp.Password
}

func (mp *MongoDBParams) GetDatabase() string {
	return mp.Database
}

func (mp *MongoDBParams) GetOptions() *options.ClientOptions {
	if mp.Options != nil {
		return mp.Options
	}
	return options.Client()
}

func (mp *MongoDBParams) GetLogger() *zap.Logger {
	return mp.Logger
}

type MongoProvider interface {
	Client() *mongo.Client
	Collection() *mongo.Collection
}
