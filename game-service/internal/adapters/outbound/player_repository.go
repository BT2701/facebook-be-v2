package outbound

import (
	"context"
	"errors"
	"game-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlayerRepository interface {
	CreatePlayer(player *models.Player) (*models.Player, error)
	GetPlayerByID(id string) (*models.Player, error)
	GetOrCreateByUserID(userID string, startingBalance float64) (*models.Player, error)
	GetPlayersByGameID(gameID string) ([]*models.Player, error)
	UpdatePlayer(player *models.Player) (*models.Player, error)
	DeletePlayer(id string) error
	GetAllPlayers() ([]*models.Player, error)
	UpdateBalance(playerID string, amount float64) (float64, error)
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
	if objectID, err := primitive.ObjectIDFromHex(id); err == nil {
		err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&player)
		if err == nil {
			return &player, nil
		}
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	err := r.collection.FindOne(ctx, bson.M{"user_id": id}).Decode(&player)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) GetOrCreateByUserID(userID string, startingBalance float64) (*models.Player, error) {
	player, err := r.GetPlayerByID(userID)
	if err == nil {
		return player, nil
	}
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	created := &models.Player{
		ID:      primitive.NewObjectID(),
		UserID:  userID,
		Balance: startingBalance,
	}
	return r.CreatePlayer(created)
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

func (r *playerRepository) GetAllPlayers() ([]*models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var players []*models.Player
	if err = cursor.All(ctx, &players); err != nil {
		return nil, err
	}

	return players, nil
}

func (r *playerRepository) UpdateBalance(playerID string, amount float64) (float64, error) {
	player, err := r.GetOrCreateByUserID(playerID, amount)
	if err != nil {
		return 0, err
	}
	player.Balance = amount
	updated, err := r.UpdatePlayer(player)
	if err != nil {
		return 0, err
	}
	return updated.Balance, nil
}
