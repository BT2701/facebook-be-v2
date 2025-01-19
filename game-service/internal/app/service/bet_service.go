package service

import (
	"game-service/internal/models"
	"game-service/internal/adapters/outbound"
)

type BetService interface {
	CreateBet(bet *models.BetOption) (*models.BetOption, error)
	GetBetByID(id string) (*models.BetOption, error)
	GetBetsByPlayerID(playerID string) ([]*models.BetOption, error)
	UpdateBet(bet *models.BetOption) (*models.BetOption, error)
	DeleteBet(id string) error
}

type betService struct {
	repo outbound.BetOptionRepository
}

func NewBetService(repo outbound.BetOptionRepository) BetService {
	return &betService{repo}
}

func (s *betService) CreateBet(bet *models.BetOption) (*models.BetOption, error) {
	return s.repo.CreateBetOption(bet)
}

func (s *betService) GetBetByID(id string) (*models.BetOption, error) {
	return s.repo.GetBetOptionByID(id)
}

func (s *betService) GetBetsByPlayerID(playerID string) ([]*models.BetOption, error) {
	return s.repo.GetBetOptionsByPlayerID(playerID)
}

func (s *betService) UpdateBet(bet *models.BetOption) (*models.BetOption, error) {
	return s.repo.UpdateBetOption(bet)
}

func (s *betService) DeleteBet(id string) error {
	return s.repo.DeleteBetOption(id)
}
