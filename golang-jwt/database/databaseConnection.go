package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect(mongoURL string) error {
	if mongoURL == "" {
		return fmt.Errorf("MONGODB_URL is not set")
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Connect(ctx); err != nil {
		return err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return err
	}
	Client = client
	fmt.Println("Connected to MongoDB!")
	return nil
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		databaseName = "jwt_auth"
	}
	return client.Database(databaseName).Collection(collectionName)
}
