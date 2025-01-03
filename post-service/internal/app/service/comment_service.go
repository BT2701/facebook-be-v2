package service

import (
	"post-service/internal/adapters/outbound"
	"post-service/internal/model"
)

type CommentService interface {
	CreateComment(comment *model.Comment) error
	GetComment(id string) (*model.Comment, error)
	UpdateComment(comment *model.Comment) error
	DeleteComment(id string) error
}

type commentService struct {
	commentRepository outbound.CommentRepository
}

func NewCommentService(commentRepository outbound.CommentRepository) CommentService {
	return &commentService{commentRepository: commentRepository}
}

func (service *commentService) CreateComment(comment *model.Comment) error {
	return service.commentRepository.CreateComment(comment)
}

func (service *commentService) GetComment(id string) (*model.Comment, error) {
	return service.commentRepository.GetComment(id)
}

func (service *commentService) UpdateComment(comment *model.Comment) error {
	return service.commentRepository.UpdateComment(comment)
}

func (service *commentService) DeleteComment(id string) error {
	return service.commentRepository.DeleteComment(id)
}

