package services

import (
	"context"
	// "errors"
	"snake_api/models"
	"snake_api/repositories"
	// "snake_api/utils"
	// "time"

	"github.com/go-redis/redis/v8"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "golang.org/x/crypto/bcrypt"
	// "github.com/golang-jwt/jwt/v5"
	"os"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type ImageService interface {
	InsertImage(ctx context.Context, Image models.Image) error
	FindAllImages(ctx context.Context) ([]models.Image, error)
	DeleteAllImages(ctx context.Context) error
	EditImage(ctx context.Context, email string, Image models.Image) error
	GetImageByUserID(ctx context.Context, id string) (*models.Image, error)
	GetImageByPostID(ctx context.Context, id string) (*models.Image, error)
	DeleteAvatar(ctx context.Context, id string) error
	DeleteImageByPostID(ctx context.Context, id string) error
}

type ImageServiceImpl struct {
	repo        repositories.ImageRepository
	redisClient *redis.Client
}

func NewImageService(repo repositories.ImageRepository, redisClient *redis.Client) ImageService {
	return &ImageServiceImpl{repo: repo, redisClient: redisClient}
}

func (s *ImageServiceImpl) InsertImage(ctx context.Context, Image models.Image) error {
	err := s.repo.InsertImage(ctx, Image)
	if err != nil {
		return err
	}
	return nil
}

func (s *ImageServiceImpl) FindAllImages(ctx context.Context) ([]models.Image, error) {
	Images, err := s.repo.FindAllImages(ctx)
	if err != nil {
		return nil, err
	}
	return Images, nil
}

func (s *ImageServiceImpl) DeleteAllImages(ctx context.Context) error {
	err := s.repo.DeleteAllImages(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *ImageServiceImpl) EditImage(ctx context.Context, email string, Image models.Image) error {
	err := s.repo.EditImage(ctx, email, Image)
	if err != nil {
		return err
	}
	return nil
}

func (s *ImageServiceImpl) GetImageByUserID(ctx context.Context, id string) (*models.Image, error) {

	Image, err := s.repo.GetImageByUserID(ctx, id)
	if err != nil {
		return nil, err
	}
	return Image, nil
}

func (s *ImageServiceImpl) GetImageByPostID(ctx context.Context, id string) (*models.Image, error) {
	
	Image, err := s.repo.GetImageByPostID(ctx, id)
	if err != nil {
		return nil, err
	}
	return Image, nil
}

func (s *ImageServiceImpl) DeleteAvatar(ctx context.Context, id string) error {
	err := s.repo.DeleteAvatar(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ImageServiceImpl) DeleteImageByPostID(ctx context.Context, id string) error {
	err := s.repo.DeleteImageByPostID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}