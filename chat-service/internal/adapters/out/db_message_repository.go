package out

import (
    "context"
    "chat-service/internal/model"
    "chat-service/pkg/database"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "log"
	"os"
)

type MongoMessageRepository struct {
    Collection *mongo.Collection
}

func NewMongoMessageRepository() *MongoMessageRepository {
    return &MongoMessageRepository{
        Collection: database.MongoClient.Database(os.Getenv("DB_NAME")).Collection("messages"),
    }
}

// SaveMessage saves a message to MongoDB
func (repo *MongoMessageRepository) SaveMessage(message *model.Message) error {
    _, err := repo.Collection.InsertOne(context.Background(), message)
    return err
}

// GetMessagesByUser fetches messages between two users
func (repo *MongoMessageRepository) GetMessagesByUser(sender, receiver string) ([]model.Message, error) {
    var messages []model.Message
    filter := bson.M{
        "$or": []bson.M{
            {"sender": sender, "receiver": receiver},
            {"sender": receiver, "receiver": sender},
        },
    }
    cursor, err := repo.Collection.Find(context.Background(), filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var message model.Message
        if err := cursor.Decode(&message); err != nil {
            log.Println("Error decoding message:", err)
            continue
        }
        messages = append(messages, message)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return messages, nil
}
func (repo *MongoMessageRepository) GetAllMessages() ([]model.Message, error) {
    var messages []model.Message
    cursor, err := repo.Collection.Find(context.Background(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var message model.Message
        if err := cursor.Decode(&message); err != nil {
            log.Println("Error decoding message:", err)
            continue
        }
        messages = append(messages, message)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return messages, nil
}
func (repo *MongoMessageRepository) DeleteAllMessages() error {
    _, err := repo.Collection.DeleteMany(context.Background(), bson.M{})
    return err
}

