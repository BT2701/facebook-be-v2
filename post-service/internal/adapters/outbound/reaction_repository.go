package outbound

import (
	"post-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"errors"
)

type ReactionRepository interface {
	CreateReaction(reaction *model.Reaction) error
	GetReaction(id string) (*model.Reaction, error)
	UpdateReaction(reaction *model.Reaction) error
	DeleteReaction(id string) error
}

type reactionRepository struct {
	collection *mongo.Collection
}

func NewReactionRepository(collection *mongo.Collection) ReactionRepository {
	return &reactionRepository{collection: collection}
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
	_, err := repo.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}