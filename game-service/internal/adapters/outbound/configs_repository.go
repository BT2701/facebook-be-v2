package outbound

import (
	"context"
	"game-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConfigsRepository interface {
	GetConfig(gameName string) (models.Common, error)
	SetConfig(gameName string, value models.Common) error
}

type configsRepository struct {
	collection *mongo.Collection
}

func NewConfigsRepository(collection *mongo.Collection) ConfigsRepository {
	return &configsRepository{
		collection: collection,
	}
}

func (r *configsRepository) GetConfig(gameName string) (models.Common, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var config models.Common
	err := r.collection.FindOne(ctx, bson.M{"game_name": gameName}).Decode(&config)
	if err != nil {
		return models.Common{}, err
	}
	return config, nil
}

func (r *configsRepository) SetConfig(gameName string, value models.Common) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"game_name": gameName}, bson.M{"$set": value})
	return err
}
