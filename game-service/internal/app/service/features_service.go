package service

import (
	"game-service/internal/adapters/outbound"
	"game-service/internal/models"
)

type FeaturesService interface {
	GetFeature(gameName string) (models.Common, error)
	SetFeature(gameName string, value models.Common) error
}

type featuresService struct {
	repo outbound.FeaturesRepository
}

func NewFeaturesService(repo outbound.FeaturesRepository) FeaturesService {
	return &featuresService{
		repo: repo,
	}
}

func (s *featuresService) GetFeature(gameName string) (models.Common, error) {
	return s.repo.GetFeature(gameName)
}

func (s *featuresService) SetFeature(gameName string, value models.Common) error {
	return s.repo.SetFeature(gameName, value)
}
