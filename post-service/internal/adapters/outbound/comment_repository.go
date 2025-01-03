package outbound

import (
	"post-service/internal/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
	"context"
	"errors"
)

type CommentRepository interface {
	CreateComment(comment *model.Comment) error
	GetComment(id string) (*model.Comment, error)
	UpdateComment(comment *model.Comment) error
	DeleteComment(id string) error
}

type commentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository(collection *mongo.Collection) CommentRepository {
	return &commentRepository{collection: collection}
}

func (repo *commentRepository) CreateComment(comment *model.Comment) error {
	_, err := repo.collection.InsertOne(context.Background(),comment)
	return err
}

func (repo *commentRepository) GetComment(id string) (*model.Comment, error) {
	var comment model.Comment
	err := repo.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (repo *commentRepository) UpdateComment(comment *model.Comment) error {
	_, err := repo.collection.ReplaceOne(context.Background(), bson.M{"_id": comment.ID}, comment)
	return err
}

func (repo *commentRepository) DeleteComment(id string) error {
	_, err := repo.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
