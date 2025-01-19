package outbound

import (
	"context"
	"time"
	"game-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BetOptionRepository interface {
	CreateBetOption(bet *models.BetOption) (*models.BetOption, error)
	GetBetOptionByID(id string) (*models.BetOption, error)
	GetBetOptionsByPlayerID(playerID string) ([]*models.BetOption, error)
	UpdateBetOption(bet *models.BetOption) (*models.BetOption, error)
	DeleteBetOption(id string) error
}

type betRepository struct {
	collection *mongo.Collection
}

func NewBetOptionRepository(collection *mongo.Collection) BetOptionRepository {
	return &betRepository{collection}
}

func (r *betRepository) CreateBetOption(bet *models.BetOption) (*models.BetOption, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, bet)
	if err != nil {
		return nil, err
	}

	return bet, nil
}

func (r *betRepository) GetBetOptionByID(id string) (*models.BetOption, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var bet models.BetOption
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bet)
	if err != nil {
		return nil, err
	}

	return &bet, nil
}

func (r *betRepository) GetBetOptionsByPlayerID(playerID string) ([]*models.BetOption, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"player_id": playerID})
	if err != nil {
		return nil, err
	}

	var bets []*models.BetOption
	if err = cursor.All(ctx, &bets); err != nil {
		return nil, err
	}

	return bets, nil
}

func (r *betRepository) UpdateBetOption(bet *models.BetOption) (*models.BetOption, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": bet.ID}, bet)
	if err != nil {
		return nil, err
	}

	return bet, nil
}

func (r *betRepository) DeleteBetOption(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

