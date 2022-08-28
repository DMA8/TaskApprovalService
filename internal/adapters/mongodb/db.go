package mongodb

import (
	"context"
	"fmt"
	"os"

	"gitlab.com/g6834/team31/tasks/internal/config"
	"gitlab.com/g6834/team31/auth/pkg/client/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	TasksCollection *mongo.Collection
}

func New(ctx context.Context, cfg config.MongoConfig) (*Database, error) {
	connStr := generateConnstr()
	clientMongo, err := mongodb.MongoClient(ctx, connStr)
	fmt.Println(connStr)
	if err != nil {
		return nil, err
	}
	collectionMongo := mongodb.MongoCollection(clientMongo, cfg.DB, cfg.TasksCollection)
	return &Database{TasksCollection: collectionMongo}, nil
}

func generateConnstr() string {
	cfg := config.NewConfig()
	mongoPass := os.Getenv("MONGO_PASSWORD")
	if mongoPass == "" {
		return cfg.Mongo.URIFul
	}
	connStr := fmt.Sprintf("mongodb://%s:%s@%s/%s", cfg.Mongo.Login, mongoPass, cfg.Mongo.URI, cfg.Mongo.DB)
	return connStr
}
