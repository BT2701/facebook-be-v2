package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("MongoDB connection error: ", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping error: ", err)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "mediadb"
	}
	DB = client.Database(dbName)
	log.Println("Connected to MongoDB")
}

func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}
