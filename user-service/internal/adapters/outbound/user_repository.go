package outbound

import (
	"context"
	"user-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"errors"
	"fmt"
)

type UserRepository interface {
	FindUserByEmailAndPassword(ctx context.Context, email, password string) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	CountUsersByEmail(ctx context.Context, email string) (int64, error)
	InsertUser(ctx context.Context, user models.User) error
	UpdateUserPassword(ctx context.Context, email, password string) error
	FindAllUsers(ctx context.Context) ([]models.User, error)
	DeleteAllUsers(ctx context.Context) error
	Logout(ctx context.Context, email string) error
	EditUser(ctx context.Context, email string, user models.User) error
    GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateAvatar(ctx context.Context, email, avatar string) error
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

func (r *userRepositoryImpl) DeleteAllUsers(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{})
	return err
}
func (r *userRepositoryImpl) Logout(ctx context.Context, email string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"is_online": 0}, "$currentDate": bson.M{"last_active": true}})
	return err
}
func (r *userRepositoryImpl) EditUser(ctx context.Context, email string, user models.User) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"name": user.Name, "description": user.Description, "birthday": user.Birthday,
		"address": user.Address, "social": user.Social, "education": user.Education, "relationship": user.Relationship, "phone": user.Phone, "gender": user.Gender, "email": user.Email}})
	return err
}
func (r *userRepositoryImpl) GetUserByID(ctx context.Context, id string) (*models.User, error) {
    fmt.Println("Received ID:", id)

    var user models.User
    filter := bson.M{"_id": id} // Không chuyển sang ObjectID
    err := r.collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            fmt.Println("No user found for ID:", id)
            return nil, errors.New("user not found")
        }
        fmt.Println("Error querying database:", err)
        return nil, err
    }

    fmt.Println("User found:", user)
    return &user, nil
}

func (r *userRepositoryImpl) UpdateAvatar(ctx context.Context, email, avatar string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"avatar": avatar}})
	return err
}