package service

import (
	"game-service/internal/adapters/outbound"
	"game-service/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type BackupService interface {
	BackupAll() error
}

type backupService struct {
	repo outbound.BackupRepository
	db   *mongo.Database
}

func NewBackupService(repo outbound.BackupRepository, db *mongo.Database) BackupService {
	return &backupService{
		repo: repo,
		db:   db,
	}
}

func (s *backupService) BackupAll() error {
	return s.repo.BackupAllCollections(models.Items, s.db)
}
