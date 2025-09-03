package service

import (
	"game-service/internal/adapters/outbound"
	"game-service/internal/models"
)

type PaylinesService interface {
	GetPayline(gameName string) (models.Common, error)
	SetPayline(gameName string, value models.Common) error
}

type paylinesService struct {
	repo outbound.PaylinesRepository
}

func NewPaylinesService(repo outbound.PaylinesRepository) PaylinesService {
	return &paylinesService{
		repo: repo,
	}
}

func (s *paylinesService) GetPayline(gameName string) (models.Common, error) {
	return s.repo.GetPayline(gameName)
}

func (s *paylinesService) SetPayline(gameName string, value models.Common) error {
	return s.repo.SetPayline(gameName, value)
}
