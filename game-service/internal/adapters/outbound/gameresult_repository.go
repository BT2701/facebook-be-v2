package outbound

import (
	"context"
	"game-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameResultRepository interface {
	CreateGameResult(game *models.GameResult) (*models.GameResult, error)
	GetGameResultByID(id string) (*models.GameResult, error)
	GetGameResultsByPlayerID(playerID string) ([]*models.GameResult, error)
	UpdateGameResult(game *models.GameResult) (*models.GameResult, error)
	DeleteGameResult(id string) error
	GetGameResultsBySessionID(sessionID string) ([]*models.GameResult, error)
}

type gameResultRepository struct {
	collection *mongo.Collection
}

func NewGameResultRepository(collection *mongo.Collection) GameResultRepository {
	return &gameResultRepository{collection}
}

func (r *gameResultRepository) CreateGameResult(game *models.GameResult) (*models.GameResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *gameResultRepository) GetGameResultByID(id string) (*models.GameResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	var game models.GameResult
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *gameResultRepository) GetGameResultsByPlayerID(playerID string) ([]*models.GameResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": game.ID}, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *gameResultRepository) DeleteGameResult(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (r *gameResultRepository) GetGameResultsBySessionID(sessionID string) ([]*models.GameResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"session_id": objectID})
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
