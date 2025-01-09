package service

import (
	"post-service/internal/adapters/outbound"
	"post-service/internal/model"
)

type StoryService interface {
	CreateStory(story *model.Story) (*model.Story, error)
	GetStory(id string) (*model.Story, error)
	UpdateStory(story *model.Story) error
	DeleteStory(id string) error
	GetStoriesByUserID(userID string) ([]model.Story, error)
	GetStories() ([]model.Story, error)
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
	return service.storyRepository.GetStoriesByUserID(userID)
}

func (service *storyService) GetStories() ([]model.Story, error) {
	return service.storyRepository.GetStories()
}
