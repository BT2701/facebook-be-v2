package service

import (
	"game-service/internal/adapters/outbound"
	"game-service/internal/models"
)

type ReelsService interface {
	GetReels(gameName string) (models.Common, error)
	SetReels(gameName string, value models.Common) error
}

type reelsService struct {
	repo outbound.ReelsRepository
}

func NewReelsService(repo outbound.ReelsRepository) ReelsService {
	return &reelsService{
		repo: repo,
	}
}

func (s *reelsService) GetReels(gameName string) (models.Common, error) {
	return s.repo.GetReel(gameName)
}

func (s *reelsService) SetReels(gameName string, value models.Common) error {
	return s.repo.SetReel(gameName, value)
}
