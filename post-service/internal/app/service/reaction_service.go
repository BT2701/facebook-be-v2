package service 

import (
	"post-service/internal/adapters/outbound"
	"post-service/internal/model"
)

type ReactionService interface {
	CreateReaction(reaction *model.Reaction) error
	GetReaction(id string) (*model.Reaction, error)
	UpdateReaction(reaction *model.Reaction) error
	DeleteReaction(id string) error
}

type reactionService struct {
	reactionRepository outbound.ReactionRepository
}

func NewReactionService(reactionRepository outbound.ReactionRepository) ReactionService {
	return &reactionService{reactionRepository: reactionRepository}
}

func (service *reactionService) CreateReaction(reaction *model.Reaction) error {
	return service.reactionRepository.CreateReaction(reaction)
}

func (service *reactionService) GetReaction(id string) (*model.Reaction, error) {
	return service.reactionRepository.GetReaction(id)
}

func (service *reactionService) UpdateReaction(reaction *model.Reaction) error {
	return service.reactionRepository.UpdateReaction(reaction)
}

func (service *reactionService) DeleteReaction(id string) error {
	return service.reactionRepository.DeleteReaction(id)
}