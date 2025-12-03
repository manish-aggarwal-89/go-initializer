package di

import (
	"context"
	"time"
	"{{MODULE_NAME}}/config"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(cfg *config.Config, logger *logrus.Logger) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	integratorDbConfig := cfg.MongoConfig.IntegratorMongoDBConfig
	clientOptions := options.Client().ApplyURI(integratorDbConfig.Uri)
	clientOptions.SetConnectTimeout(integratorDbConfig.ConnectionTimeout)
	clientOptions.SetMaxPoolSize(uint64(integratorDbConfig.MaxPool))
	clientOptions.SetMinPoolSize(uint64(integratorDbConfig.MinPool))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal(err)
	}

	return client.Database(integratorDbConfig.DBName), nil
}
