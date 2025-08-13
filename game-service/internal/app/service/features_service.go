package service

import (
	"fmt"
	configsDb "game-service/internal/adapters/db"
)

type FeaturesService interface {
	GetFeature(key string) (string, error)
	SetFeature(key string, value string) error
	GetAllFeatures() (map[string]string, error)
}

func NewFeaturesService() FeaturesService {
	return &featuresService{
		features: make(map[string]string),
	}
}

type featuresService struct {
	features map[string]string
}

func (s *featuresService) GetFeature(key string) (string, error) {
	value, ok := s.features[key]
	if !ok {
		return "", fmt.Errorf("feature key %s not found", key)
	}
	return value, nil
}

func (s *featuresService) SetFeature(key string, value string) error {
	s.features[key] = value
	return nil
}

func (s *featuresService) GetAllFeatures() (map[string]string, error) {
	features, err := configsDb.NewMongoDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get all features: %w", err)
	}
	return features, nil
}