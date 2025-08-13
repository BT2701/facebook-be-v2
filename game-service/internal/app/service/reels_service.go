package service

import (
	"fmt"
	configsDb "game-service/internal/adapters/db"
)

type ReelsService interface {
	GetReel(key string) (string, error)
	SetReel(key string, value string) error
	GetAllReels() (map[string]string, error)
}

func NewReelsService() ReelsService {
	return &reelsService{
		reels: make(map[string]string),
	}
}

type reelsService struct {
	reels map[string]string
}

func (s *reelsService) GetReel(key string) (string, error) {
	value, ok := s.reels[key]
	if !ok {
		return "", fmt.Errorf("reel key %s not found", key)
	}
	return value, nil
}

func (s *reelsService) SetReel(key string, value string) error {
	s.reels[key] = value
	return nil
}

func (s *reelsService) GetAllReels() (map[string]string, error) {
	reels, err := configsDb.NewMongoDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get all reels: %w", err)
	}
	return reels, nil
}
