package service

import (
	"game-service/internal/adapters/outbound"
	"game-service/internal/models"
)

type ConfigsService interface {
	GetConfig(gameName string) (models.Common, error)
	SetConfig(gameName string, value models.Common) error
}

type configsService struct {
	repo outbound.ConfigsRepository
}

func NewConfigsService(repo outbound.ConfigsRepository) ConfigsService {
	return &configsService{
		repo: repo,
	}
}

func (s *configsService) GetConfig(gameName string) (models.Common, error) {
	return s.repo.GetConfig(gameName)
}

func (s *configsService) SetConfig(gameName string, value models.Common) error {
	return s.repo.SetConfig(gameName, value)
}
