package outbound

import (
	"ai-service/internal/model"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationRepository interface {
	GetOrCreate(ctx context.Context, userID string) (*model.Conversation, error)
	GetByID(ctx context.Context, id, userID string) (*model.Conversation, error)
	ListByUser(ctx context.Context, userID string) ([]model.Conversation, error)
	Touch(ctx context.Context, id string) error
	Delete(ctx context.Context, id, userID string) error
}

type conversationRepository struct {
	collection *mongo.Collection
}

func NewConversationRepository(collection *mongo.Collection) ConversationRepository {
	repo := &conversationRepository{collection: collection}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_, _ = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return repo
}

func (r *conversationRepository) GetOrCreate(ctx context.Context, userID string) (*model.Conversation, error) {
	var existing model.Conversation
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&existing)
	if err == nil {
		return &existing, nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	now := time.Now()
	conversation := model.Conversation{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Title:     "AI assistant",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := r.collection.InsertOne(ctx, conversation); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			if findErr := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&existing); findErr == nil {
				return &existing, nil
			}
		}
		return nil, err
	}
	return &conversation, nil
}

func (r *conversationRepository) GetByID(ctx context.Context, id, userID string) (*model.Conversation, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var conversation model.Conversation
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID, "user_id": userID}).Decode(&conversation)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *conversationRepository) ListByUser(ctx context.Context, userID string) ([]model.Conversation, error) {
	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversations []model.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, err
	}
	if conversations == nil {
		conversations = []model.Conversation{}
	}
	return conversations, nil
}

func (r *conversationRepository) Touch(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateByID(ctx, objectID, bson.M{"$set": bson.M{"updated_at": time.Now()}})
	return err
}

func (r *conversationRepository) Delete(ctx context.Context, id, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID, "user_id": userID})
	return err
}
