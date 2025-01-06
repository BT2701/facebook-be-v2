package database

import (
    "context"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var MongoClient *mongo.Client

// InitMongoDB initializes MongoDB connection
func InitMongoDB() *mongo.Client {
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }

    log.Println("Successfully connected to MongoDB!")
    MongoClient = client
    return MongoClient
}

func GetCollection(databaseName, collectionName string) *mongo.Collection {
    return MongoClient.Database(databaseName).Collection(collectionName)
}
