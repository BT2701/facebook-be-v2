package outbound

import (
	"post-service/internal/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
	"context"
	"errors"
)

type StoryRepository interface {
	CreateStory(story *model.Story) (*model.Story, error)
	GetStory(id string) (*model.Story, error)
	UpdateStory(story *model.Story) error
	DeleteStory(id string) error
	GetStoriesByUserID(userID string) ([]model.Story, error)
	GetStories() ([]model.Story, error)
	DeleteStories () error
}

type storyRepository struct {
	collection *mongo.Collection
}

func NewStoryRepository(collection *mongo.Collection) StoryRepository {
	return &storyRepository{collection: collection}
}

func (repo *storyRepository) CreateStory(story *model.Story) (*model.Story, error) {
	_, err := repo.collection.InsertOne(context.Background(), story)
	if err != nil {
		return nil, err
	}
	return story, nil
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

func (repo *storyRepository) GetStoriesByUserID(userID string) ([]model.Story, error) {
	var stories []model.Story
	cursor, err := repo.collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &stories); err != nil {
		return nil, err
	}
	return stories, nil
}

func (repo *storyRepository) GetStories() ([]model.Story, error) {
	var stories []model.Story
	cursor, err := repo.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &stories); err != nil {
		return nil, err
	}
	return stories, nil
}

func (repo *storyRepository) DeleteStories() error {
	_, err := repo.collection.DeleteMany(context.Background(), bson.M{})
	return err
}
