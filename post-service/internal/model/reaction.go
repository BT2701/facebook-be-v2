package model
import "time"
import "go.mongodb.org/mongo-driver/bson/primitive"

type Reaction struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    UserID    string             `json:"user_id" bson:"user_id"`
    Timeline  time.Time          `json:"timeline" bson:"timeline"`
    PostID    primitive.ObjectID `json:"post_id" bson:"post_id,omitempty"` // Khóa ngoại trỏ tới Post
    CommentID primitive.ObjectID `json:"comment_id" bson:"comment_id,omitempty"` // Khóa ngoại trỏ tới Comment (nếu cần)
}

