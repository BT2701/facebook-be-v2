package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Title     string             `json:"title" bson:"title"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Message struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ConversationID string             `json:"conversation_id" bson:"conversation_id"`
	Role           string             `json:"role" bson:"role"`
	Content        string             `json:"content" bson:"content"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}
