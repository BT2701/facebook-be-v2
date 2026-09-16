package outbound

import (
	"post-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"errors"
	"time"
)

type ReactionRepository interface {
	CreateReaction(reaction *model.Reaction) error
	GetReaction(id string) (*model.Reaction, error)
	UpdateReaction(reaction *model.Reaction) error
	DeleteReaction(id string) error
	GetByPostAndUser(postID, userID string) (*model.Reaction, error)
}

type reactionRepository struct {
	collection *mongo.Collection
}

func NewReactionRepository(collection *mongo.Collection) ReactionRepository {
	repo := &reactionRepository{collection: collection}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_, _ = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "post_id", Value: 1}, {Key: "user_id", Value: 1}},
		Options: options.Index(),
	})
	return repo
}

func (repo *reactionRepository) CreateReaction(reaction *model.Reaction) error {
	_, err := repo.collection.InsertOne(context.Background(),reaction)
	return err
}

func (repo *reactionRepository) GetReaction(id string) (*model.Reaction, error) {
	var reaction model.Reaction
	err := repo.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&reaction)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &reaction, nil
}

func (repo *reactionRepository) UpdateReaction(reaction *model.Reaction) error {
	_, err := repo.collection.ReplaceOne(context.Background(), bson.M{"_id": reaction.ID}, reaction)
	return err
}

func (repo *reactionRepository) DeleteReaction(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = repo.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

func (repo *reactionRepository) GetByPostAndUser(postID, userID string) (*model.Reaction, error) {
	objectID, err := primitive.ObjectIDFromHex(postID)
	filter := bson.M{"user_id": userID, "post_id": postID}
	if err == nil {
		filter = bson.M{"user_id": userID, "$or": []bson.M{{"post_id": objectID}, {"post_id": postID}}}
	}
	var reaction model.Reaction
	err = repo.collection.FindOne(context.Background(), filter).Decode(&reaction)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &reaction, nil
}