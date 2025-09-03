package outbound

import (
	"context"
	"game-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SymbolsRepository interface {
	GetSymbol(gameName string) (models.Common, error)
	SetSymbol(gameName string, value models.Common) error
}

type symbolsRepository struct {
	collection *mongo.Collection
}

func NewSymbolsRepository(collection *mongo.Collection) SymbolsRepository {
	return &symbolsRepository{
		collection: collection,
	}
}

func (r *symbolsRepository) GetSymbol(gameName string) (models.Common, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	var symbol models.Common
	err := r.collection.FindOne(ctx, bson.M{"game_name": gameName}).Decode(&symbol)
	if err != nil {
		return models.Common{}, err
	}
	return symbol, nil
}

func (r *symbolsRepository) SetSymbol(gameName string, value models.Common) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"game_name": gameName}, bson.M{"$set": value})
	return err
}
