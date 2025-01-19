package service

import (
	"game-service/internal/models"
	"game-service/internal/adapters/outbound"
)

type SymbolService interface {
	CreateSymbol(symbol *models.Symbol) (*models.Symbol, error)
	GetSymbolByID(id string) (*models.Symbol, error)
	GetSymbolsByGameID(gameID string) ([]*models.Symbol, error)
	UpdateSymbol(symbol *models.Symbol) (*models.Symbol, error)
	DeleteSymbol(id string) error
}

type symbolService struct {
	repo outbound.SymbolRepository
}

func NewSymbolService(repo outbound.SymbolRepository) SymbolService {
	return &symbolService{repo}
}

func (s *symbolService) CreateSymbol(symbol *models.Symbol) (*models.Symbol, error) {
	return s.repo.CreateSymbol(symbol)
}

func (s *symbolService) GetSymbolByID(id string) (*models.Symbol, error) {
	return s.repo.GetSymbolByID(id)
}

func (s *symbolService) GetSymbolsByGameID(gameID string) ([]*models.Symbol, error) {
	return s.repo.GetSymbolsByGameID(gameID)
}

func (s *symbolService) UpdateSymbol(symbol *models.Symbol) (*models.Symbol, error) {
	return s.repo.UpdateSymbol(symbol)
}

func (s *symbolService) DeleteSymbol(id string) error {
	return s.repo.DeleteSymbol(id)
}
