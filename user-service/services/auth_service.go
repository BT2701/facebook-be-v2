package services

import (
	"context"
	"errors"
	"snake_api/models"
	"snake_api/repositories"
	"snake_api/utils"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (string, error)
	SignUp(ctx context.Context, user models.User) error
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type userServiceImpl struct {
	repo        repositories.UserRepository
	redisClient *redis.Client
}

func NewUserService(repo repositories.UserRepository, redisClient *redis.Client) UserService {
	return &userServiceImpl{repo: repo, redisClient: redisClient}
}

func (s *userServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	// Kiểm tra token đã tồn tại trong Redis
	existingToken, err := s.redisClient.Get(ctx, "user:"+email).Result()
	if err == nil && existingToken != "" {
		return existingToken, nil // Trả về token cũ
	}

	// Lấy thông tin người dùng từ MongoDB
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Xác thực mật khẩu
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Tạo token mới
	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	// Lưu token vào Redis với TTL (e.g., 24 giờ)
	err = s.redisClient.Set(ctx, "user:"+email, token, 24*time.Hour).Err()
	if err != nil {
		return "", errors.New("failed to save token in Redis")
	}

	return token, nil
}

func (s *userServiceImpl) SignUp(ctx context.Context, user models.User) error {
	// Kiểm tra email đã tồn tại
	count, _ := s.repo.CountUsersByEmail(ctx, user.Email)
	if count > 0 {
		return errors.New("user already exists")
	}

	// Băm mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	// Thêm thông tin khác
	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = time.Now()

	// Lưu người dùng vào MongoDB
	return s.repo.InsertUser(ctx, user)
}

func (s *userServiceImpl) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}

	resetToken, err := utils.GenerateTokenWithExpiry(user.Email, time.Hour*1)
	if err != nil {
		return "", errors.New("failed to generate reset token")
	}

	resetURL := "http://localhost:3000/reset-password?token=" + resetToken
	err = utils.SendEmail(email, "Password Reset Request", "Click here to reset your password: "+resetURL)
	if err != nil {
		return "", errors.New("failed to send email")
	}

	return resetToken, nil
}

func (s *userServiceImpl) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Giải mã token để lấy email
	claims, err := utils.DecodeToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	email := claims["email"].(string)

	// Băm mật khẩu mới
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	// Cập nhật mật khẩu trong MongoDB
	return s.repo.UpdateUserPassword(ctx, email, string(hashedPassword))
}
func (s *userServiceImpl) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.FindAllUsers(ctx)
}
