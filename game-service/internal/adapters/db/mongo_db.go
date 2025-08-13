package db

import (
	"context"
	dbModels "game-service/internal/models"
	"game-service/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	DB          *mongo.Database
	GameDbData  dbModels.GameModel
	BackupItems []dbModels.BackupItem
}

func NewMongoDB(db *mongo.Database) *MongoDB {
	return &MongoDB{
		DB:          db,
		BackupItems: dbModels.Items,
	}
}

func (m *MongoDB) LoadReelsDB() {
	var reels dbModels.Common
	collection := m.DB.Collection("reels")
	err := collection.FindOne(context.Background(), bson.M{}).Decode(&reels)
	if err == nil {
		m.GameDbData.Reels = reels
	} else{
		m.BackupReelsDB()
	}
}

func (m *MongoDB) LoadPaylinesDB() {
	var paylines dbModels.Common
	collection := m.DB.Collection("paylines")
	err := collection.FindOne(context.Background(), bson.M{}).Decode(&paylines)
	if err == nil {
		m.GameDbData.Paylines = paylines
	} else {
		m.BackupPaylinesDB()
	}
}

func (m *MongoDB) LoadSymbolsDB() {
	var symbols dbModels.Common
	collection := m.DB.Collection("symbols")
	err := collection.FindOne(context.Background(), bson.M{}).Decode(&symbols)
	if err == nil {
		m.GameDbData.Symbols = symbols
	} else {
		m.BackupSymbolsDB()
	}
}

func (m *MongoDB) LoadConfigsDB() {
	var configs dbModels.Common
	collection := m.DB.Collection("configs")
	err := collection.FindOne(context.Background(), bson.M{}).Decode(&configs)
	if err == nil {
		m.GameDbData.Configs = configs
	} else {
		m.BackupConfigsDB()
	}
}

func (m *MongoDB) LoadFeaturesDB() {
	var features dbModels.Common
	collection := m.DB.Collection("features")
	err := collection.FindOne(context.Background(), bson.M{}).Decode(&features)
	if err == nil {
		m.GameDbData.Features = features
	} else {
		m.BackupFeaturesDB()
	}
}

func (m *MongoDB) BackupReelsDB() error {
	var data dbModels.Common
	err := utils.LoadJSONData(m.BackupItems[1].FilePath, &data)
	if err != nil {
		return err
	}
	coll := m.DB.Collection(m.BackupItems[1].Collection)
	// Xóa dữ liệu cũ
	_, err = coll.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	// Thêm dữ liệu mới
	_, err = coll.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) BackupPaylinesDB() error {
	var data dbModels.Common
	err := utils.LoadJSONData(m.BackupItems[2].FilePath, &data)
	if err != nil {
		return err
	}
	coll := m.DB.Collection(m.BackupItems[2].Collection)
	// Xóa dữ liệu cũ
	_, err = coll.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	// Thêm dữ liệu mới
	_, err = coll.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) BackupSymbolsDB() error {
	var data dbModels.Common
	err := utils.LoadJSONData(m.BackupItems[3].FilePath, &data)
	if err != nil {
		return err
	}
	coll := m.DB.Collection(m.BackupItems[3].Collection)
	// Xóa dữ liệu cũ
	_, err = coll.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	// Thêm dữ liệu mới
	_, err = coll.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) BackupConfigsDB() error {
	var data dbModels.Common
	err := utils.LoadJSONData(m.BackupItems[0].FilePath, &data)
	if err != nil {
		return err
	}
	coll := m.DB.Collection(m.BackupItems[0].Collection)
	// Xóa dữ liệu cũ
	_, err = coll.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	// Thêm dữ liệu mới
	_, err = coll.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) BackupFeaturesDB() error {
	var data dbModels.Common
	err := utils.LoadJSONData(m.BackupItems[4].FilePath, &data)
	if err != nil {
		return err
	}
	coll := m.DB.Collection(m.BackupItems[4].Collection)
	// Xóa dữ liệu cũ
	_, err = coll.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	// Thêm dữ liệu mới
	_, err = coll.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}
	return nil
}
