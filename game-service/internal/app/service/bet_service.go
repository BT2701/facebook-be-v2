package service

import (
	"game-service/internal/models"
	"game-service/pkg/utils"
	"fmt"
)

const betOptionsFilePath = "pkg/json/bet_options.json"

type BetService struct{}

func NewBetService() *BetService {
	return &BetService{}
}

func (s *BetService) LoadBetOptions() (*models.GameBetOptions, error) {
	var gameBetOptionData struct {
		GameName string                `json:"game_name"`
		Data     models.GameBetOptions `json:"data"`
	}

	err := utils.LoadJSONData(betOptionsFilePath, &gameBetOptionData)
	if err != nil {
		return nil, fmt.Errorf("failed to load bet options: %w", err)
	}

	return &gameBetOptionData.Data, nil
}
