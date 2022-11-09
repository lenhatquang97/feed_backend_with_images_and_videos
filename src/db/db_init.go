package db

import (
	"context"
	"log"
	"time"

	"example.com/feed_backend/src/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDb *mongo.Client = ConnectDB()
var feedCollection *mongo.Collection = GetCollection(mongoDb, "feed")

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("feed_database").Collection(collectionName)
	return collection
}

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(configs.EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
