package outbound

import (
	"context"
	"errors"
	"post-service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository interface {
	CreatePost(post *model.Post) error
	GetPost(id string) (*model.Post, error)
	UpdatePost(post *model.Post) error
	DeletePost(id string) error
	GetPostsByUserID(userID string) ([]model.Post, error)
	GetPosts() ([]model.Post, error)
	DeleteAllPosts() error
	SearchPosts(content string, limit, offset int64) ([]model.Post, error)
}

type postRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(collection *mongo.Collection) PostRepository {
	repo := &postRepository{collection: collection}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_, _ = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "timeline", Value: -1}}},
		{Keys: bson.D{{Key: "timeline", Value: -1}}},
	})
	return repo
}

func (repo *postRepository) CreatePost(post *model.Post) error {
	_, err := repo.collection.InsertOne(context.Background(), post)
	return err
}

func (repo *postRepository) GetPost(id string) (*model.Post, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var post model.Post
	err = repo.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&post)
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
	opts := options.Find().SetSort(bson.D{{Key: "timeline", Value: -1}}).SetLimit(100)
	cursor, err := repo.collection.Find(context.Background(), bson.M{}, opts)
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

func (repo *postRepository) DeleteAllPosts() error {
	_, err := repo.collection.DeleteMany(context.Background(), bson.M{})
	return err
}

func (repo *postRepository) SearchPosts(content string, limit, offset int64) ([]model.Post, error) {
	if limit <= 0 {
		limit = 20
	}
	filter := bson.M{}
	if content != "" {
		filter["content"] = bson.M{"$regex": content, "$options": "i"}
	}
	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"timeline": -1})
	cursor, err := repo.collection.Find(context.Background(), filter, opts)
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
