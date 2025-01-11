package outbound

import (
	"context"
	"media-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ImageRepository interface {
	InsertImage(ctx context.Context, Image models.Image) error
	FindAllImages(ctx context.Context) ([]models.Image, error)
	DeleteAllImages(ctx context.Context) error
	EditImage(ctx context.Context, email string, Image models.Image) error
    GetImageByUserID(ctx context.Context, id string) (*models.Image, error)
	GetImageByPostID(ctx context.Context, id string) (*models.Image, error)
	DeleteAvatar(ctx context.Context, id string) error
	DeleteImageByPostID(ctx context.Context, id string) error
}

type ImageRepositoryImpl struct {
	collection *mongo.Collection
}

func NewImageRepository(collection *mongo.Collection) ImageRepository {
	return &ImageRepositoryImpl{collection: collection}
}

func (r *ImageRepositoryImpl) InsertImage(ctx context.Context, Image models.Image) error {
	_, err := r.collection.InsertOne(ctx, Image)
	if err != nil {
		return err
	}
	return nil
}

func (r *ImageRepositoryImpl) FindAllImages(ctx context.Context) ([]models.Image, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var Images []models.Image

	for cursor.Next(ctx) {
		var Image models.Image
		if err := cursor.Decode(&Image); err != nil {
			return nil, err
		}
		Images = append(Images, Image)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return Images, nil
}

func (r *ImageRepositoryImpl) DeleteAllImages(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}
	return nil
}

func (r *ImageRepositoryImpl) EditImage(ctx context.Context, email string, Image models.Image) error {
	_, err := r.collection.UpdateOne(ctx , bson.M{"email": email}, bson.M{"$set": Image})
	if err != nil {
		return err
	}
	return nil
}

func (r *ImageRepositoryImpl) GetImageByUserID(ctx context.Context, id string) (*models.Image, error) {
	var Image models.Image
	err := r.collection.FindOne(ctx, bson.M{"user_id": id}).Decode(&Image)
	if err != nil {
		return nil, err
	}
	return &Image, nil
}

func (r *ImageRepositoryImpl) GetImageByPostID(ctx context.Context, id string) (*models.Image, error) {
	var Image models.Image
	err := r.collection.FindOne(ctx, bson.M{"post_id": id}).Decode(&Image)
	if err != nil {
		return nil, err
	}
	return &Image, nil
}

func (r *ImageRepositoryImpl) DeleteAvatar(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (r *ImageRepositoryImpl) DeleteImageByPostID(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"post_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (r *ImageRepositoryImpl) FindImageByPostID(ctx context.Context, id string) (*models.Image, error) {
	var Image models.Image
	err := r.collection.FindOne(ctx, bson.M{"post_id": id}).Decode(&Image)
	if err != nil {
		return nil, err
	}
	return &Image, nil
}

func (r *ImageRepositoryImpl) FindImageByUserID(ctx context.Context, id string) (*models.Image, error) {
	var Image models.Image
	err := r.collection.FindOne(ctx, bson.M{"user_id": id}).Decode(&Image)
	if err != nil {
		return nil, err
	}
	return &Image, nil
}


func (r *ImageRepositoryImpl) DeleteImageByUserID(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": id})
	if err != nil {
		return err
	}
	return nil
}

