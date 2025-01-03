package outbound

import (
	"post-service/internal/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
	"context"
	"errors"
)

type StoryRepository interface {
	CreateStory(story *model.Story) error
	GetStory(id string) (*model.Story, error)
	UpdateStory(story *model.Story) error
	DeleteStory(id string) error
}

type storyRepository struct {
	collection *mongo.Collection
}

func NewStoryRepository(collection *mongo.Collection) StoryRepository {
	return &storyRepository{collection: collection}
}

func (repo *storyRepository) CreateStory(story *model.Story) error {
	_, err := repo.collection.InsertOne(context.Background(),story)
	return err
}

func (repo *storyRepository) GetStory(id string) (*model.Story, error) {
	var story model.Story
	err := repo.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&story)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &story, nil
}

func (repo *storyRepository) UpdateStory(story *model.Story) error {
	_, err := repo.collection.ReplaceOne(context.Background(), bson.M{"_id": story.ID}, story)
	return err
}

func (repo *storyRepository) DeleteStory(id string) error {
	_, err := repo.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

