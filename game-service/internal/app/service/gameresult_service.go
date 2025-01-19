package service

import (
	"game-service/internal/models"
	"game-service/internal/adapters/outbound"
)

type GameResultService interface {
	CreateGameResult(game *models.GameResult) (*models.GameResult, error)
	GetGameResultByID(id string) (*models.GameResult, error)
	GetGameResultsByPlayerID(playerID string) ([]*models.GameResult, error)
	UpdateGameResult(game *models.GameResult) (*models.GameResult, error)
	DeleteGameResult(id string) error
}

type gameResultService struct {
	repo outbound.GameResultRepository
}

func NewGameResultService(repo outbound.GameResultRepository) GameResultService {
	return &gameResultService{repo}
}

func (s *gameResultService) CreateGameResult(game *models.GameResult) (*models.GameResult, error) {
	return s.repo.CreateGameResult(game)
}

func (s *gameResultService) GetGameResultByID(id string) (*models.GameResult, error) {
	return s.repo.GetGameResultByID(id)
}

func (s *gameResultService) GetGameResultsByPlayerID(playerID string) ([]*models.GameResult, error) {
	return s.repo.GetGameResultsByPlayerID(playerID)
}

func (s *gameResultService) UpdateGameResult(game *models.GameResult) (*models.GameResult, error) {
	return s.repo.UpdateGameResult(game)
}

func (s *gameResultService) DeleteGameResult(id string) error {
	return s.repo.DeleteGameResult(id)
}

