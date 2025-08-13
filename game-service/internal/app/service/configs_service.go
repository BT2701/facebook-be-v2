package service

import (
	"fmt"
	configsDb "game-service/internal/adapters/db"
)

type ConfigsService interface {
	GetConfig(key string) (string, error)
	SetConfig(key string, value string) error
	GetAllConfigs() (map[string]string, error)
}

func NewConfigsService() ConfigsService {
	return &configsService{
		configs: make(map[string]string),
	}
}

type configsService struct {
	configs map[string]string
}

func (s *configsService) GetConfig(key string) (string, error) {
	value, ok := s.configs[key]
	if !ok {
		return "", fmt.Errorf("config key %s not found", key)
	}
	return value, nil
}

func (s *configsService) SetConfig(key string, value string) error {
	s.configs[key] = value
	return nil
}

func (s *configsService) GetAllConfigs() (map[string]string, error) {
	configs, err := configsDb.NewMongoDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get all configs: %w", err)
	}
	return configs, nil
}