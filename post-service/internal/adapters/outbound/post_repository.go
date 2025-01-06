package outbound

import (
	"post-service/internal/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
	"context"
	"errors"
)

type PostRepository interface {
	CreatePost(post *model.Post) error
	GetPost(id string) (*model.Post, error)
	UpdatePost(post *model.Post) error
	DeletePost(id string) error
	GetPostsByUserID(userID string) ([]model.Post, error)
	GetPosts() ([]model.Post, error)
}

type postRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(collection *mongo.Collection) PostRepository {
	return &postRepository{collection: collection}
}

func (repo *postRepository) CreatePost(post *model.Post) error {
	_, err := repo.collection.InsertOne(context.Background(), post)
	return err
}

func (repo *postRepository) GetPost(id string) (*model.Post, error) {
	var post model.Post
	err := repo.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (repo *postRepository) UpdatePost(post *model.Post) error {
	_, err := repo.collection.ReplaceOne(context.Background(), bson.M{"_id": post.ID}, post)
	return err
}

func (repo *postRepository) DeletePost(id string) error {
	_, err := repo.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (repo *postRepository) GetPostsByUserID(userID string) ([]model.Post, error) {
	cursor, err := repo.collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var posts []model.Post
	if err := cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *postRepository) GetPosts() ([]model.Post, error) {
	cursor, err := repo.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var posts []model.Post
	if err := cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
}