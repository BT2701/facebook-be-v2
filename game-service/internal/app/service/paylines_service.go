package service

import (
	"fmt"
	configsDb "game-service/internal/adapters/db"
)

type PaylinesService interface {
	GetPayline(key string) (string, error)
	SetPayline(key string, value string) error
	GetAllPaylines() (map[string]string, error)
}

func NewPaylinesService() PaylinesService {
	return &paylinesService{
		paylines: make(map[string]string),
	}
}

type paylinesService struct {
	paylines map[string]string
}

func (s *paylinesService) GetPayline(key string) (string, error) {
	value, ok := s.paylines[key]
	if !ok {
		return "", fmt.Errorf("payline key %s not found", key)
	}
	return value, nil
}

func (s *paylinesService) SetPayline(key string, value string) error {
	s.paylines[key] = value
	return nil
}

func (s *paylinesService) GetAllPaylines() (map[string]string, error) {
	paylines, err := configsDb.NewMongoDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get all paylines: %w", err)
	}
	return paylines, nil
}
