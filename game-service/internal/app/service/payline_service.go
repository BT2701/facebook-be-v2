package service

import (
	"fmt"
	"game-service/internal/models"
	"game-service/pkg/utils"
)

const paylineFilePath = "pkg/json/payline.json"

type PaylineService struct {
	symbolService *SymbolService
}

func NewPaylineService() *PaylineService {
	return &PaylineService{}
}

type GamePaylineData struct {
	GameName string              `json:"game_name" bson:"game_name"`
	Data     models.GamePaylines `json:"data" bson:"data"`
}

func (s *PaylineService) LoadPaylines() (*models.GamePaylines, error) {
	var gamePaylineData GamePaylineData

	// Dùng LoadJSONData từ utils để load file
	err := utils.LoadJSONData(paylineFilePath, &gamePaylineData)
	if err != nil {
		return nil, fmt.Errorf("failed to load paylines: %w", err)
	}

	return &gamePaylineData.Data, nil
}

func (s *PaylineService) CalculateWinnings(finalSymbols [][]string, betAmount float64) (float64, []map[string]interface{}) {
	var winAmount float64
	var results []map[string]interface{}
	paylinesData, err := s.LoadPaylines()
	if err != nil {
		fmt.Println("Error loading paylines:", err)
		return 0, nil
	}

	winningPatterns := []int{5, 4, 3}

	for _, payline := range paylinesData.Paylines {
		pattern := payline.Pattern
		maxCols := len(pattern)
		if payline.Direction == "left_to_right" {
			for _, size := range winningPatterns {
				if size > maxCols {
					continue
				}

				subsetPattern := pattern[:size]
				initialSymbol := finalSymbols[subsetPattern[0][1]][subsetPattern[0][0]]

				if initialSymbol == "" {
					continue
				}

				isWinningPattern := true
				for _, pos := range subsetPattern {
					col, row := pos[0], pos[1]
					if finalSymbols[row][col] != initialSymbol {
						isWinningPattern = false
						break
					}
				}

				if isWinningPattern {
					symbolValue := s.symbolService.GetSymbolValue(initialSymbol, size)
					winAmount += betAmount * symbolValue

					results = append(results, map[string]interface{}{
						"paylineId":   payline.ID,
						"symbol":      initialSymbol,
						"occurrences": size,
						"onWinline":   true,
						"symbolValue": symbolValue,
						"positions":   subsetPattern,
					})
					break
				}
			}
		}
	}

	return winAmount, results
}
