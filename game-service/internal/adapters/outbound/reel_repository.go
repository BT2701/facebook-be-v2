package outbound

import (
	"context"
	"time"
	"game-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReelRepository interface {
	CreateReel(reel *models.Reel) (*models.Reel, error)
	GetReelByID(id string) (*models.Reel, error)
	GetReelsByGameID(gameID string) ([]*models.Reel, error)
	UpdateReel(reel *models.Reel) (*models.Reel, error)
	DeleteReel(id string) error
}

type reelRepository struct {
	collection *mongo.Collection
}

func NewReelRepository(collection *mongo.Collection) ReelRepository {
	return &reelRepository{collection}
}

func (r *reelRepository) CreateReel(reel *models.Reel) (*models.Reel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, reel)
	if err != nil {
		return nil, err
	}

	return reel, nil
}

func (r *reelRepository) GetReelByID(id string) (*models.Reel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var reel models.Reel
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&reel)
	if err != nil {
		return nil, err
	}

	return &reel, nil
}

func (r *reelRepository) GetReelsByGameID(gameID string) ([]*models.Reel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"game_id": gameID})
	if err != nil {
		return nil, err
	}

	var reels []*models.Reel
	if err = cursor.All(ctx, &reels); err != nil {
		return nil, err
	}

	return reels, nil
}

func (r *reelRepository) UpdateReel(reel *models.Reel) (*models.Reel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": reel.ID}, reel)
	if err != nil {
		return nil, err
	}

	return reel, nil
}

func (r *reelRepository) DeleteReel(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
