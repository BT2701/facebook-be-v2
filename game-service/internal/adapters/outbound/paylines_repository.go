package outbound

import (
	"context"
	"game-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaylinesRepository interface {
	GetPayline(gameName string) (models.Common, error)
	SetPayline(gameName string, value models.Common) error
}

type paylinesRepository struct {
	collection *mongo.Collection
}

func NewPaylinesRepository(collection *mongo.Collection) PaylinesRepository {
	return &paylinesRepository{
		collection: collection,
	}
}

func (r *paylinesRepository) GetPayline(gameName string) (models.Common, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	var payline models.Common
	err := r.collection.FindOne(ctx, bson.M{"game_name": gameName}).Decode(&payline)
	if err != nil {
		return models.Common{}, err
	}
	return payline, nil
}

func (r *paylinesRepository) SetPayline(gameName string, value models.Common) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"game_name": gameName}, bson.M{"$set": value})
	return err
}
