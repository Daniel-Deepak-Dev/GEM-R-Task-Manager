package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Could not connect to MongoDB: ", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Could not ping MongoDB: ", err)
	}

	log.Println("Connected to MongoDB successfully")
	return client
}
