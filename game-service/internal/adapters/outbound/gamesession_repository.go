package outbound

import (
	"context"
	"time"
	"game-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type GameSessionRepository interface {
	CreateGameSession(game *models.GameSession) (*models.GameSession, error)
	GetGameSessionByID(id string) (*models.GameSession, error)
	GetGameSessionsByPlayerID(playerID string, page, limit int) ([]*models.GameSession, error)
	UpdateGameSession(game *models.GameSession) (*models.GameSession, error)
	DeleteGameSession(id string) error
}

type gameSessionRepository struct {
	collection *mongo.Collection
}

func NewGameSessionRepository(collection *mongo.Collection) GameSessionRepository {
	return &gameSessionRepository{collection}
}

func (r *gameSessionRepository) CreateGameSession(game *models.GameSession) (*models.GameSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *gameSessionRepository) GetGameSessionByID(id string) (*models.GameSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var game models.GameSession
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *gameSessionRepository) GetGameSessionsByPlayerID(playerID string, page, limit int) ([]*models.GameSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(playerID)
	if err != nil {
		return nil, err
	}

	skip := (page - 1) * limit
	opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"player_id": objectID}, opts)
	if err != nil {
		return nil, err
	}

	var games []*models.GameSession
	err = cursor.All(ctx, &games)
	if err != nil {
		return nil, err
	}

	return games, nil
}


func (r *gameSessionRepository) UpdateGameSession(game *models.GameSession) (*models.GameSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": game.ID}, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *gameSessionRepository) DeleteGameSession(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

