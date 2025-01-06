package service

import (
	"post-service/internal/adapters/outbound"
	"post-service/internal/model"
	"time"
)

type PostService interface {
	CreatePost(post *model.Post) error
	GetPost(id string) (*model.Post, error)
	UpdatePost(post *model.Post) error
	DeletePost(id string) error
	GetPostsByUserID(userID string) ([]model.Post, error)
	GetPosts() ([]model.Post, error)
}

type postService struct {
	postRepository outbound.PostRepository
}

func NewPostService(postRepository outbound.PostRepository) PostService {
	return &postService{postRepository: postRepository}
}

func (service *postService) CreatePost(post *model.Post) error {
	post.Timeline = time.Now()
	return service.postRepository.CreatePost(post)
}

func (service *postService) GetPost(id string) (*model.Post, error) {
	return service.postRepository.GetPost(id)
}

func (service *postService) UpdatePost(post *model.Post) error {
	return service.postRepository.UpdatePost(post)
}

func (service *postService) DeletePost(id string) error {
	return service.postRepository.DeletePost(id)
}

func (service *postService) GetPostsByUserID(userID string) ([]model.Post, error) {
	return service.postRepository.GetPostsByUserID(userID)
}

func (service *postService) GetPosts() ([]model.Post, error) {
	return service.postRepository.GetPosts()
}
