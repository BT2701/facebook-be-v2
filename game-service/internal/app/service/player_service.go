package service

import (
	"game-service/internal/models"
	"game-service/internal/adapters/outbound"
)

type PlayerService interface {
	CreatePlayer(player *models.Player) (*models.Player, error)
	GetPlayerByID(id string) (*models.Player, error)
	GetPlayersByGameID(gameID string) ([]*models.Player, error)
	UpdatePlayer(player *models.Player) (*models.Player, error)
	DeletePlayer(id string) error
}

type playerService struct {
	repo outbound.PlayerRepository
}

func NewPlayerService(repo outbound.PlayerRepository) PlayerService {
	return &playerService{repo}
}

func (s *playerService) CreatePlayer(player *models.Player) (*models.Player, error) {
	return s.repo.CreatePlayer(player)
}

func (s *playerService) GetPlayerByID(id string) (*models.Player, error) {
	return s.repo.GetPlayerByID(id)
}

func (s *playerService) GetPlayersByGameID(gameID string) ([]*models.Player, error) {
	return s.repo.GetPlayersByGameID(gameID)
}

func (s *playerService) UpdatePlayer(player *models.Player) (*models.Player, error) {
	return s.repo.UpdatePlayer(player)
}

func (s *playerService) DeletePlayer(id string) error {
	return s.repo.DeletePlayer(id)
}

