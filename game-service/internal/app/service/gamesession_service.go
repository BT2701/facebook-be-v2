package service

import (
	"game-service/internal/models"
	"game-service/internal/adapters/outbound"
)

type GameSessionService interface {
	CreateGameSession(game *models.GameSession) (*models.GameSession, error)
	GetGameSessionByID(id string) (*models.GameSession, error)
	GetGameSessionsByPlayerID(playerID string, page, limit int) ([]*models.GameSession, error)
	UpdateGameSession(game *models.GameSession) (*models.GameSession, error)
	DeleteGameSession(id string) error
}

type gameSessionService struct {
	repo outbound.GameSessionRepository
}

func NewGameSessionService(repo outbound.GameSessionRepository) GameSessionService {
	return &gameSessionService{repo}
}

func (s *gameSessionService) CreateGameSession(game *models.GameSession) (*models.GameSession, error) {
	return s.repo.CreateGameSession(game)
}

func (s *gameSessionService) GetGameSessionByID(id string) (*models.GameSession, error) {
	return s.repo.GetGameSessionByID(id)
}

func (s *gameSessionService) GetGameSessionsByPlayerID(playerID string, page, limit int) ([]*models.GameSession, error) {
	return s.repo.GetGameSessionsByPlayerID(playerID, page, limit)
}


func (s *gameSessionService) UpdateGameSession(game *models.GameSession) (*models.GameSession, error) {
	return s.repo.UpdateGameSession(game)
}

func (s *gameSessionService) DeleteGameSession(id string) error {
	return s.repo.DeleteGameSession(id)
}
