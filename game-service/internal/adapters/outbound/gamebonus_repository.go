package outbound

import (
	"context"
	"time"
	"game-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BonusGameRepository interface {
	CreateBonusGame(game *models.BonusGame) (*models.BonusGame, error)
	GetBonusGameByID(id string) (*models.BonusGame, error)
	GetBonusGamesByPlayerID(playerID string) ([]*models.BonusGame, error)
	UpdateBonusGame(game *models.BonusGame) (*models.BonusGame, error)
	DeleteBonusGame(id string) error
}

type bonusGameRepository struct {
	collection *mongo.Collection
}

func NewBonusGameRepository(collection *mongo.Collection) BonusGameRepository {
	return &bonusGameRepository{collection}
}

func (r *bonusGameRepository) CreateBonusGame(game *models.BonusGame) (*models.BonusGame, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *bonusGameRepository) GetBonusGameByID(id string) (*models.BonusGame, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var game models.BonusGame
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *bonusGameRepository) GetBonusGamesByPlayerID(playerID string) ([]*models.BonusGame, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"player_id": playerID})
	if err != nil {
		return nil, err
	}

	var games []*models.BonusGame
	if err = cursor.All(ctx, &games); err != nil {
		return nil, err
	}

	return games, nil
}

func (r *bonusGameRepository) UpdateBonusGame(game *models.BonusGame) (*models.BonusGame, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": game.ID}, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *bonusGameRepository) DeleteBonusGame(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}


