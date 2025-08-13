package service

import (
	"fmt"
	configsDb "game-service/internal/adapters/db"
)

type SymbolsService interface {
	GetSymbol(key string) (string, error)
	SetSymbol(key string, value string) error
	GetAllSymbols() (map[string]string, error)
}

func NewSymbolsService() SymbolsService {
	return &symbolsService{
		symbols: make(map[string]string),
	}
}

type symbolsService struct {
	symbols map[string]string
}

func (s *symbolsService) GetSymbol(key string) (string, error) {
	value, ok := s.symbols[key]
	if !ok {
		return "", fmt.Errorf("symbol key %s not found", key)
	}
	return value, nil
}

func (s *symbolsService) SetSymbol(key string, value string) error {
	s.symbols[key] = value
	return nil
}

func (s *symbolsService) GetAllSymbols() (map[string]string, error) {
	symbols, err := configsDb.NewMongoDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get all symbols: %w", err)
	}
	return symbols, nil
}
