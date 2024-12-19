package repositories

import (
	"context"
	"snake_api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindUserByEmailAndPassword(ctx context.Context, email, password string) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	CountUsersByEmail(ctx context.Context, email string) (int64, error)
	InsertUser(ctx context.Context, user models.User) error
	UpdateUserPassword(ctx context.Context, email, password string) error
}

type userRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepositoryImpl{collection: collection}
}

func (r *userRepositoryImpl) FindUserByEmailAndPassword(ctx context.Context, email, password string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user)
	return &user, err
}

func (r *userRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (r *userRepositoryImpl) CountUsersByEmail(ctx context.Context, email string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"email": email})
}

func (r *userRepositoryImpl) InsertUser(ctx context.Context, user models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepositoryImpl) UpdateUserPassword(ctx context.Context, email, password string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"password": password}})
	return err
}
