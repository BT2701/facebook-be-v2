package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type GameModel struct {
	Configs   Common      `bson:"configs" json:"configs"`
	Reels     Common      `bson:"reels" json:"reels"`
	Paylines  Common      `bson:"paylines" json:"paylines"`
	Symbols   Common      `bson:"symbols" json:"symbols"`
	Histories interface{} `bson:"histories" json:"histories"`
	Users     interface{} `bson:"users" json:"users"`
	Features  Common      `bson:"features" json:"features"`
}

type Common struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GameName string             `bson:"game_name" json:"game_name"`
	Data     interface{}        `bson:"data" json:"data"`
}

var Items = []BackupItem{
	{"game-service/backup_db/configs.json", "configs"},
	{"game-service/backup_db/reels.json", "reels"},
	{"game-service/backup_db/paylines.json", "paylines"},
	{"game-service/backup_db/symbols.json", "symbols"},
	{"game-service/backup_db/features.json", "features"},
}

type BackupItem struct {
	FilePath   string
	Collection string
}