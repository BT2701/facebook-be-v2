package service

import (
	"fmt"
	"game-service/internal/models"
	"game-service/pkg/utils"
)

const symbolsFilePath = "pkg/json/symbols.json"

type SymbolService struct{}

func NewSymbolService() *SymbolService {
	return &SymbolService{}
}

type GameSymbols struct {
	GameName string              `json:"game_name"`
	Data     models.GameSymbolData `json:"data"`
}

func (s *SymbolService)LoadSymbols() (*models.GameSymbolData, error) {
	var gameSymbols GameSymbols

	err := utils.LoadJSONData(symbolsFilePath, &gameSymbols)
	if err != nil {
		return nil, fmt.Errorf("failed to load symbols: %w", err)
	}

	return &gameSymbols.Data, nil
}

func (s *SymbolService) GetSymbolValue(symbolName string, count int) float64 {
	symbolData, err := s.LoadSymbols()
	if err != nil {
		return 0.0 // Trường hợp không load được symbols
	}
	
	for _, symbol := range symbolData.Symbols {
		if symbol.Name == symbolName {
			return symbol.Values[fmt.Sprint(count)]
		}
	}
	return 0.0
}

