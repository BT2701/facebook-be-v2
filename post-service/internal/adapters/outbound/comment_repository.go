package outbound

import (
	"post-service/internal/model"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"errors"
	"time"
)

type CommentRepository interface {
	CreateComment(comment *model.Comment) error
	GetComment(id string) (*model.Comment, error)
	UpdateComment(comment *model.Comment) error
	DeleteComment(id string) error
	GetCommentsByPostID(postID string) ([]model.Comment, error)
}

type commentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository(collection *mongo.Collection) CommentRepository {
	repo := &commentRepository{collection: collection}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_, _ = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "post_id", Value: 1}},
		Options: options.Index(),
	})
	return repo
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

func (repo *commentRepository) GetCommentsByPostID(postID string) ([]model.Comment, error) {
	objectID, err := primitive.ObjectIDFromHex(postID)
	filter := bson.M{"post_id": postID}
	if err == nil {
		filter = bson.M{"$or": []bson.M{{"post_id": objectID}, {"post_id": postID}}}
	}
	cursor, err := repo.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var comments []model.Comment
	if err := cursor.All(context.Background(), &comments); err != nil {
		return nil, err
	}
	return comments, nil
}
