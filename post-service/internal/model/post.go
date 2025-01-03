package model

import "time"
import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
    ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`  // Không có dấu cách dư sau "omitempty"
    Content   string               `json:"content" bson:"content"`
    Timeline  time.Time            `json:"timeline" bson:"timeline"`
    UserID    string               `json:"user_id" bson:"user_id"`
    Image     string               `json:"image" bson:"image"`
    Comments  []primitive.ObjectID `json:"comments" bson:"comments,omitempty"`   // Mảng các ObjectID của Comment
    Reactions []primitive.ObjectID `json:"reactions" bson:"reactions,omitempty"` // Mảng các ObjectID của Reaction
}

