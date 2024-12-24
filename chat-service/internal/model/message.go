package model

import "time"

// Message represents a chat message.
type Message struct {
    ID        string    `json:"id" bson:"_id"`
    Sender    string    `json:"sender" bson:"sender"`
    Receiver  string    `json:"receiver" bson:"receiver"`
    Content   string    `json:"content" bson:"content"`
    Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
