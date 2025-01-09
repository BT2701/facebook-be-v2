package model

import "time"

// Message represents a chat message.
type Message struct {
    ID        string    `json:"id" bson:"_id"`
    Sender    string    `json:"sender" bson:"sender"`
    Receiver  string    `json:"receiver" bson:"receiver"`
    Content   string    `json:"content" bson:"content"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
