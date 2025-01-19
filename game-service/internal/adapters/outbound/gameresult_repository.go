package outbound

import (
	"context"
	"time"
	"game-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameResultRepository interface {
	CreateGameResult(game *models.GameResult) (*models.GameResult, error)
	GetGameResultByID(id string) (*models.GameResult, error)
	GetGameResultsByPlayerID(playerID string) ([]*models.GameResult, error)
	UpdateGameResult(game *models.GameResult) (*models.GameResult, error)
	DeleteGameResult(id string) error
}

type gameResultRepository struct {
	collection *mongo.Collection
}

func NewGameResultRepository(collection *mongo.Collection) GameResultRepository {
	return &gameResultRepository{collection}
}

func (r *gameResultRepository) CreateGameResult(game *models.GameResult) (*models.GameResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *gameResultRepository) GetGameResultByID(id string) (*models.GameResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var game models.GameResult
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *gameResultRepository) GetGameResultsByPlayerID(playerID string) ([]*models.GameResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"player_id": playerID})
	if err != nil {
		return nil, err
	}

	var games []*models.GameResult
	for cursor.Next(ctx) {
		var game models.GameResult
		err := cursor.Decode(&game)
		if err != nil {
			return nil, err
		}
		games = append(games, &game)
	}

	return games, nil
}

func (r *gameResultRepository) UpdateGameResult(game *models.GameResult) (*models.GameResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": game.ID}, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *gameResultRepository) DeleteGameResult(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
