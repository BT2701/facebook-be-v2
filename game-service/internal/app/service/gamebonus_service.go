package service

import (
	"game-service/internal/models"
	"game-service/internal/adapters/outbound"
)

type BonusGameService interface {
	CreateBonusGame(bonusGame *models.BonusGame) (*models.BonusGame, error)
	GetBonusGameByID(id string) (*models.BonusGame, error)
	GetBonusGamesByPlayerID(playerID string) ([]*models.BonusGame, error)
	UpdateBonusGame(bonusGame *models.BonusGame) (*models.BonusGame, error)
	DeleteBonusGame(id string) error
}

type bonusGameService struct {
	repo outbound.BonusGameRepository
}

func NewBonusGameService(repo outbound.BonusGameRepository) BonusGameService {
	return &bonusGameService{repo}
}

func (s *bonusGameService) CreateBonusGame(bonusGame *models.BonusGame) (*models.BonusGame, error) {
	return s.repo.CreateBonusGame(bonusGame)
}

func (s *bonusGameService) GetBonusGameByID(id string) (*models.BonusGame, error) {
	return s.repo.GetBonusGameByID(id)
}

func (s *bonusGameService) GetBonusGamesByPlayerID(playerID string) ([]*models.BonusGame, error) {
	return s.repo.GetBonusGamesByPlayerID(playerID)
}

func (s *bonusGameService) UpdateBonusGame(bonusGame *models.BonusGame) (*models.BonusGame, error) {
	return s.repo.UpdateBonusGame(bonusGame)
}

func (s *bonusGameService) DeleteBonusGame(id string) error {
	return s.repo.DeleteBonusGame(id)
}
