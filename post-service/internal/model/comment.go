package model

import "time"
import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
    ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
    UserID    string               `json:"user_id" bson:"user_id"`
    Content   string               `json:"content" bson:"content"`
    Timeline  time.Time            `json:"timeline" bson:"timeline"`
    PostID    primitive.ObjectID   `json:"post_id" bson:"post_id,omitempty"`  // Có thể bỏ dấu cách dư sau "omitempty"
    Reactions []primitive.ObjectID `json:"reactions" bson:"reactions,omitempty"`
}

