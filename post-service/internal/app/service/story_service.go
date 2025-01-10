package service

import (
	"post-service/internal/adapters/outbound"
	"post-service/internal/model"
	"time"
)

type StoryService interface {
	CreateStory(story *model.Story) (*model.Story, error)
	GetStory(id string) (*model.Story, error)
	UpdateStory(story *model.Story) error
	DeleteStory(id string) error
	GetStoriesByUserID(userID string) ([]model.Story, error)
	GetStories() ([]model.Story, error)
	DeleteStories() error
}

type storyService struct {
	storyRepository outbound.StoryRepository
}

func NewStoryService(storyRepository outbound.StoryRepository) StoryService {
	return &storyService{storyRepository: storyRepository}
}

func (service *storyService) CreateStory(story *model.Story) (*model.Story, error) {
	return service.storyRepository.CreateStory(story)
}

func (service *storyService) GetStory(id string) (*model.Story, error) {
	return service.storyRepository.GetStory(id)
}

func (service *storyService) UpdateStory(story *model.Story) error {
	return service.storyRepository.UpdateStory(story)
}

func (service *storyService) DeleteStory(id string) error {
	return service.storyRepository.DeleteStory(id)
}

func (service *storyService) GetStoriesByUserID(userID string) ([]model.Story, error) {
	stories, err := service.storyRepository.GetStoriesByUserID(userID)
	if err != nil {
		return nil, err
	}

	var recentStories []model.Story
	now := time.Now()
	for _, story := range stories {
		if now.Sub(story.Timeline) <= 24*time.Hour {
			recentStories = append(recentStories, story)
		}
	}

	return recentStories, nil
}

func (service *storyService) GetStories() ([]model.Story, error) {
	return service.storyRepository.GetStories()
}

func (service *storyService) DeleteStories() error {
	return service.storyRepository.DeleteStories()
}