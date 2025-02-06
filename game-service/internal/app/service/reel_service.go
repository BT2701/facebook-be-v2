package service

import (
	"game-service/internal/models"
	"game-service/pkg/utils"
	"fmt"
)

const reelFilePath = "pkg/json/reel.json"

type GameReelData struct {
	GameName string      `json:"game_name"`
	Data     models.Reel `json:"data"`
}

type ReelService struct{}

func NewReelService() *ReelService {
	return &ReelService{}
}

func (s *ReelService) LoadReel() (*models.Reel, error) {
	var gameReelData GameReelData

	err := utils.LoadJSONData(reelFilePath, &gameReelData)
	if err != nil {
		return nil, fmt.Errorf("failed to load reel data: %w", err)
	}

	return &gameReelData.Data, nil
}
