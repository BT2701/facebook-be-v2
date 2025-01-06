package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Friend struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	UserID1  string             `json:"userID1" bson:"userID1"`
	UserID2  string             `json:"userID2" bson:"userID2"`
	IsFriend bool               `json:"isFriend" bson:"isFriend"`
	Timeline time.Time          `json:"timeline" bson:"timeline"`
}
