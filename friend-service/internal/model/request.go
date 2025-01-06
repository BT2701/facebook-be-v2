package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Request struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Sender   string             `json:"sender" bson:"sender"`
	Receiver string             `json:"receiver" bson:"receiver"`
	Timeline time.Time          `json:"timeline" bson:"timeline"`
}
