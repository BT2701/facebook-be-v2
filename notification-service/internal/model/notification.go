package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Notification represents a notification.
type Notification struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	User     string             `json:"user" bson:"user"`
	Content  string             `json:"content" bson:"content"`
	Timeline time.Time          `json:"timeline" bson:"timeline"`
	Receiver string             `json:"receiver" bson:"receiver"`
	Post     string             `json:"post" bson:"post"`
	Is_read  int                `json:"is_read" bson:"is_read"`
	Action_n int                `json:"action_n" bson:"action_n"`
}
