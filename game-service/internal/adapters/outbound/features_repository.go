package outbound

import (
	"context"
	"game-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeaturesRepository interface {
	GetFeature(gameName string) (models.Common, error)
	SetFeature(gameName string, value models.Common) error
}

type featuresRepository struct {
	collection *mongo.Collection
}

func NewFeaturesRepository(collection *mongo.Collection) FeaturesRepository {
	return &featuresRepository{
		collection: collection,
	}
}

func (r *featuresRepository) GetFeature(gameName string) (models.Common, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	var feature models.Common
	err := r.collection.FindOne(ctx, bson.M{"game_name": gameName}).Decode(&feature)
	if err != nil {
		return models.Common{}, err
	}
	return feature, nil
}

func (r *featuresRepository) SetFeature(gameName string, value models.Common) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"game_name": gameName}, bson.M{"$set": value})
	return err
}
