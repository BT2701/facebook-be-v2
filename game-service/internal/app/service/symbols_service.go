package service

import (
	"game-service/internal/adapters/outbound"
	"game-service/internal/models"
)

type SymbolsService interface {
	GetSymbol(gameName string) (models.Common, error)
	SetSymbol(gameName string, value models.Common) error
}

type symbolsService struct {
	repo outbound.SymbolsRepository
}

func NewSymbolsService(repo outbound.SymbolsRepository) SymbolsService {
	return &symbolsService{
		repo: repo,
	}
}

func (s *symbolsService) GetSymbol(gameName string) (models.Common, error) {
	return s.repo.GetSymbol(gameName)
}

func (s *symbolsService) SetSymbol(gameName string, value models.Common) error {
	return s.repo.SetSymbol(gameName, value)
}
