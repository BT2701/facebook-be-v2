package services

import (
	"context"
	"errors"
	"snake_api/models"
	"snake_api/repositories"
	"snake_api/utils"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (string, error)
	SignUp(ctx context.Context, user models.User) error
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type userServiceImpl struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindUserByEmailAndPassword(ctx, email, password)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateToken(user.Email)
}

func (s *userServiceImpl) SignUp(ctx context.Context, user models.User) error {
	count, _ := s.repo.CountUsersByEmail(ctx, user.Email)
	if count > 0 {
		return errors.New("user already exists")
	}

	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = time.Now()

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
	claims, err := utils.DecodeToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	email := claims["email"].(string)
	return s.repo.UpdateUserPassword(ctx, email, newPassword)
}
