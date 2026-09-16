package outbound

import (
	"ai-service/internal/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	ListByConversation(ctx context.Context, conversationID string, limit int64) ([]model.Message, error)
	DeleteByConversation(ctx context.Context, conversationID string) error
}

type messageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(collection *mongo.Collection) MessageRepository {
	repo := &messageRepository{collection: collection}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_, _ = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "conversation_id", Value: 1}, {Key: "created_at", Value: 1}},
	})
	return repo
}

func (r *messageRepository) Create(ctx context.Context, message *model.Message) error {
	message.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, message)
	return err
}

func (r *messageRepository) ListByConversation(ctx context.Context, conversationID string, limit int64) ([]model.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}).SetLimit(limit)
	cursor, err := r.collection.Find(ctx, bson.M{"conversation_id": conversationID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []model.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	if messages == nil {
		messages = []model.Message{}
	}
	return messages, nil
}

func (r *messageRepository) DeleteByConversation(ctx context.Context, conversationID string) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"conversation_id": conversationID})
	return err
}
