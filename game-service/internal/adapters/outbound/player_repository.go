package outbound

import (
	"context"
	"time"
	"game-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlayerRepository interface {
	CreatePlayer(player *models.Player) (*models.Player, error)
	GetPlayerByID(id string) (*models.Player, error)
	GetPlayersByGameID(gameID string) ([]*models.Player, error)
	UpdatePlayer(player *models.Player) (*models.Player, error)
	DeletePlayer(id string) error
}

type playerRepository struct {
	collection *mongo.Collection
}

func NewPlayerRepository(collection *mongo.Collection) PlayerRepository {
	return &playerRepository{collection}
}

func (r *playerRepository) CreatePlayer(player *models.Player) (*models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (r *playerRepository) GetPlayerByID(id string) (*models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var player models.Player
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&player)
	if err != nil {
		return nil, err
	}

	return &player, nil
}

func (r *playerRepository) GetPlayersByGameID(gameID string) ([]*models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"game_id": gameID})
	if err != nil {
		return nil, err
	}

	var players []*models.Player
	if err = cursor.All(ctx, &players); err != nil {
		return nil, err
	}

	return players, nil
}

func (r *playerRepository) UpdatePlayer(player *models.Player) (*models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": player.ID}, player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (r *playerRepository) DeletePlayer(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
