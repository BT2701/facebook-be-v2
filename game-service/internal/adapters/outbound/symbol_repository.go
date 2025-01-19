package outbound

import (
	"context"
	"time"
	"game-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SymbolRepository interface {
	CreateSymbol(symbol *models.Symbol) (*models.Symbol, error)
	GetSymbolByID(id string) (*models.Symbol, error)
	GetSymbolsByGameID(gameID string) ([]*models.Symbol, error)
	UpdateSymbol(symbol *models.Symbol) (*models.Symbol, error)
	DeleteSymbol(id string) error
}

type symbolRepository struct {
	collection *mongo.Collection
}

func NewSymbolRepository(collection *mongo.Collection) SymbolRepository {
	return &symbolRepository{collection}
}

func (r *symbolRepository) CreateSymbol(symbol *models.Symbol) (*models.Symbol, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, symbol)
	if err != nil {
		return nil, err
	}

	return symbol, nil
}

func (r *symbolRepository) GetSymbolByID(id string) (*models.Symbol, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var symbol models.Symbol
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&symbol)
	if err != nil {
		return nil, err
	}

	return &symbol, nil
}

func (r *symbolRepository) GetSymbolsByGameID(gameID string) ([]*models.Symbol, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"game_id": gameID})
	if err != nil {
		return nil, err
	}

	var symbols []*models.Symbol
	err = cursor.All(ctx, &symbols)
	if err != nil {
		return nil, err
	}

	return symbols, nil
}

func (r *symbolRepository) UpdateSymbol(symbol *models.Symbol) (*models.Symbol, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": symbol.ID}, symbol)
	if err != nil {
		return nil, err
	}

	return symbol, nil
}

func (r *symbolRepository) DeleteSymbol(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

