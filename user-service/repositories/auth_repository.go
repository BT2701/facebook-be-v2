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
	FindAllUsers(ctx context.Context) ([]models.User, error)
}

type userRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepositoryImpl{collection: collection}
}

func (r *userRepositoryImpl) FindAllUsers(ctx context.Context) ([]models.User, error) {
	// Tìm tất cả người dùng với bộ lọc rỗng
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err // Trả về lỗi nếu không thể thực hiện truy vấn
	}
	defer cursor.Close(ctx) // Đảm bảo đóng cursor sau khi xử lý xong

	// Khởi tạo slice để chứa kết quả
	var users []models.User

	// Lặp qua các tài liệu trong cursor và giải mã vào slice
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err // Trả về lỗi nếu giải mã thất bại
		}
		users = append(users, user)
	}

	// Kiểm tra lỗi sau khi lặp
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil // Trả về danh sách người dùng
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
