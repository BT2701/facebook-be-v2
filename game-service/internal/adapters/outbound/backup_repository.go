package outbound

import (
	"context"
	"encoding/json"
	"game-service/internal/models"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type BackupRepository interface {
	BackupAllCollections(items []models.BackupItem, db *mongo.Database) error
}

type backupRepository struct{}

func NewBackupRepository() BackupRepository {
	return &backupRepository{}
}

func (r *backupRepository) BackupAllCollections(items []models.BackupItem, db *mongo.Database) error {
	ctx := context.Background()
	for _, item := range items {
		// Kiểm tra collection có tồn tại chưa
		exists := false
		colls, err := db.ListCollectionNames(ctx, map[string]interface{}{})
		if err != nil {
			return err
		}
		for _, coll := range colls {
			if coll == item.Collection {
				exists = true
				break
			}
		}
		if !exists {
			// Tạo collection mới nếu chưa tồn tại
			err := db.CreateCollection(ctx, item.Collection)
			if err != nil {
				return err
			}
		}

		collection := db.Collection(item.Collection)

		// Đọc dữ liệu từ file JSON (object)
		file, err := os.Open(item.FilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		var doc map[string]interface{}
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&doc); err != nil {
			return err
		}

		// Insert dữ liệu vào Mongo
		_, err = collection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}
	}
	return nil
}
