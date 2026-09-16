package services

import (
	"context"
	"errors"
	"time"
	"user-service/internal/models"
	"user-service/internal/adapters/outbound"
	"user-service/pkg/utils"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (string, error, *models.User)
	SignUp(ctx context.Context, user models.User) error
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
	GetAllUsers(ctx context.Context, exclude []string, limit int64) ([]models.User, error)
	DeleteAllUsers(ctx context.Context) error
	Logout(ctx context.Context, email string) error
	EditUser(ctx context.Context, email string, user models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateAvatar(ctx context.Context, email, avatar string) error
	SearchUsers(ctx context.Context, name string, limit, offset int64) ([]models.User, error)
}

type userServiceImpl struct {
	repo        outbound.UserRepository
	redisClient *redis.Client
}

func NewUserService(repo outbound.UserRepository, redisClient *redis.Client) UserService {
	return &userServiceImpl{repo: repo, redisClient: redisClient}
}

func (s *userServiceImpl) Login(ctx context.Context, email, password string) (string, error, *models.User) {
	// Lấy thông tin người dùng từ MongoDB
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials"), nil
	}

	// Kiểm tra token đã tồn tại trong Redis
	existingToken, err := s.redisClient.Get(ctx, "user:"+email).Result()
	if err == nil && existingToken != "" {
		if claims, parseErr := utils.DecodeToken(existingToken); parseErr == nil {
			if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).After(time.Now()) {
				user.Sanitize()
				return existingToken, nil, user
			}
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials"), nil
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		return "", errors.New("failed to generate token"), nil
	}

	// Lưu token vào Redis với TTL (e.g., 24 giờ)
	err = s.redisClient.Set(ctx, "user:"+email, token, 24*time.Hour).Err()
	if err != nil {
		return "", errors.New("failed to save token in Redis"), nil
	}

	_ = s.repo.SetOnline(ctx, email, 1)
	user.Sanitize()
	return token, nil, user
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

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	resetURL := frontendURL + "/reset-password?token=" + resetToken
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
func (s *userServiceImpl) GetAllUsers(ctx context.Context, exclude []string, limit int64) ([]models.User, error) {
	users, err := s.repo.FindUsersExcluding(ctx, exclude, limit)
	if err != nil {
		return nil, err
	}
	return models.SanitizeUsers(users), nil
}

func (s *userServiceImpl) DeleteAllUsers(ctx context.Context) error {
	return s.repo.DeleteAllUsers(ctx)
}

func (s *userServiceImpl) Logout(ctx context.Context, email string) error {
	return s.repo.Logout(ctx, email)
}

func (s *userServiceImpl) EditUser(ctx context.Context, email string, user models.User) error {
	return s.repo.EditUser(ctx, email, user)
}

func (s *userServiceImpl) GetByID(ctx context.Context, id string) (*models.User, error) {
    user, err := s.repo.GetUserByID(ctx, id)
    if err != nil {
        return nil, err
    }
	user.Sanitize()
    return user, nil
}
func (s *userServiceImpl) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	user.Sanitize()
	return user, nil
}

func (s *userServiceImpl) UpdateAvatar(ctx context.Context, email, avatar string) error {
	return s.repo.UpdateAvatar(ctx, email, avatar)
}

func (s *userServiceImpl) SearchUsers(ctx context.Context, name string, limit, offset int64) ([]models.User, error) {
	users, err := s.repo.SearchUsers(ctx, name, limit, offset)
	if err != nil {
		return nil, err
	}
	return models.SanitizeUsers(users), nil
}