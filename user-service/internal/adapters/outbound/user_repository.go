package outbound

import (
	"context"
	"errors"
	"fmt"
	"time"
	"user-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	FindUserByEmailAndPassword(ctx context.Context, email, password string) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	CountUsersByEmail(ctx context.Context, email string) (int64, error)
	InsertUser(ctx context.Context, user models.User) error
	UpdateUserPassword(ctx context.Context, email, password string) error
	FindAllUsers(ctx context.Context) ([]models.User, error)
	FindUsersExcluding(ctx context.Context, exclude []string, limit int64) ([]models.User, error)
	DeleteAllUsers(ctx context.Context) error
	Logout(ctx context.Context, email string) error
	EditUser(ctx context.Context, email string, user models.User) error
    GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateAvatar(ctx context.Context, email, avatar string) error
	SearchUsers(ctx context.Context, name string, limit, offset int64) ([]models.User, error)
	SetOnline(ctx context.Context, email string, online int) error
}

type userRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	repo := &userRepositoryImpl{collection: collection}
	repo.ensureIndexes()
	return repo
}

func (r *userRepositoryImpl) ensureIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_, _ = r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "name", Value: 1}}},
	})
}

func (r *userRepositoryImpl) FindUsersExcluding(ctx context.Context, exclude []string, limit int64) ([]models.User, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	filter := bson.M{}
	if len(exclude) > 0 {
		filter["_id"] = bson.M{"$nin": exclude}
	}
	opts := options.Find().SetLimit(limit).SetProjection(bson.M{"password": 0})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
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

func (r *userRepositoryImpl) SearchUsers(ctx context.Context, name string, limit, offset int64) ([]models.User, error) {
	if limit <= 0 {
		limit = 20
	}
	filter := bson.M{}
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	opts := options.Find().SetLimit(limit).SetSkip(offset)
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepositoryImpl) SetOnline(ctx context.Context, email string, online int) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"email": email}, bson.M{
		"$set":         bson.M{"is_online": online},
		"$currentDate": bson.M{"last_active": true},
	})
	return err
}