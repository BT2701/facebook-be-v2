package service

import (
	"errors"
	"friend-service/internal/adapters/outbound"
	"friend-service/internal/model"
	"time"
)

type RequestService interface {
	CreateRequest(sender, receiver string) (*model.Request, error)
	GetRequest(sender, receiver string) (*model.Request, error)
	GetRequests(receiver string) ([]*model.Request, error)
	UpdateRequest(sender, receiver string, isAccepted bool) (*model.Request, error)
	DeleteRequest(sender, receiver string) error
}

type requestService struct {
	requestRepository outbound.RequestRepository
}

func NewRequestService(requestRepository outbound.RequestRepository) RequestService {
	return &requestService{requestRepository}
}

func (s *requestService) CreateRequest(sender, receiver string) (*model.Request, error) {
	request, err := s.requestRepository.GetRequest(sender, receiver)
	if err == nil {
		return nil, errors.New("request already exists")
	}

	request = &model.Request{
		Sender:   sender,
		Receiver: receiver,
		Timeline: time.Now(),
	}

	request, err = s.requestRepository.CreateRequest(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (s *requestService) GetRequest(sender, receiver string) (*model.Request, error) {
	request, err := s.requestRepository.GetRequest(sender, receiver)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (s *requestService) GetRequests(receiver string) ([]*model.Request, error) {
	requests, err := s.requestRepository.GetRequests(receiver)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (s *requestService) UpdateRequest(sender, receiver string, isAccepted bool) (*model.Request, error) {
	request, err := s.requestRepository.GetRequest(sender, receiver)
	if err != nil {
		return nil, err
	}

	request.Timeline = time.Now()

	request, err = s.requestRepository.UpdateRequest(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}


func (s *requestService) DeleteRequest(sender, receiver string) error {
	err := s.requestRepository.DeleteRequest(sender, receiver)
	if err != nil {
		return err
	}

	return nil
}

