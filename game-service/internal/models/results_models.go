package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type GameResult struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GameID   primitive.ObjectID `bson:"game_id" json:"game_id"`
	PlayerID primitive.ObjectID `bson:"player_id" json:"player_id"`
	Score    int                `bson:"score" json:"score"`
}
