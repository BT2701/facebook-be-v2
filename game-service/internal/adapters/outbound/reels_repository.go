package outbound

import (
	"context"
	"game-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReelsRepository interface {
	GetReel(gameName string) (models.Common, error)
	SetReel(gameName string, value models.Common) error
}

type reelsRepository struct {
	collection *mongo.Collection
}

func NewReelsRepository(collection *mongo.Collection) ReelsRepository {
	return &reelsRepository{
		collection: collection,
	}
}

func (r *reelsRepository) GetReel(gameName string) (models.Common, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var reel models.Common
	err := r.collection.FindOne(ctx, bson.M{"game_name": gameName}).Decode(&reel)
	if err != nil {
		return models.Common{}, err
	}
	return reel, nil
}

func (r *reelsRepository) SetReel(gameName string, value models.Common) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"game_name": gameName}, bson.M{"$set": value})
	return err
}
