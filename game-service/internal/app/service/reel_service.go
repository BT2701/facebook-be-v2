package service

import (
	"game-service/internal/models"
	"game-service/internal/adapters/outbound"
)

type ReelService interface {
	CreateReel(reel *models.Reel) (*models.Reel, error)
	GetReelByID(id string) (*models.Reel, error)
	GetReelsByGameID(gameID string) ([]*models.Reel, error)
	UpdateReel(reel *models.Reel) (*models.Reel, error)
	DeleteReel(id string) error
}

type reelService struct {
	repo outbound.ReelRepository
}

func NewReelService(repo outbound.ReelRepository) ReelService {
	return &reelService{repo}
}

func (s *reelService) CreateReel(reel *models.Reel) (*models.Reel, error) {
	return s.repo.CreateReel(reel)
}

func (s *reelService) GetReelByID(id string) (*models.Reel, error) {
	return s.repo.GetReelByID(id)
}

func (s *reelService) GetReelsByGameID(gameID string) ([]*models.Reel, error) {
	return s.repo.GetReelsByGameID(gameID)
}

func (s *reelService) UpdateReel(reel *models.Reel) (*models.Reel, error) {
	return s.repo.UpdateReel(reel)
}

func (s *reelService) DeleteReel(id string) error {
	return s.repo.DeleteReel(id)
}

