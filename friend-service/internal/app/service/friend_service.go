package service

import (
	"errors"
	"friend-service/internal/adapters/outbound"
	"friend-service/internal/model"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FriendService interface {
	CreateFriend(userID1, userID2 string) (*model.Friend, error)
	GetFriend(userID1, userID2 string) (*model.Friend, error)
	GetFriends(userID string) ([]*model.Friend, error)
	UpdateFriend(userID1, userID2 string, isFriend bool) (*model.Friend, error)
	DeleteFriend(userID1, userID2 string) error
	GetFriendsByUserID(userID string) ([]*model.Friend, error)
	IsFriend(userID1, userID2 string) (bool, error)
}

type friendService struct {
	friendRepository outbound.FriendRepository
}

func NewFriendService(friendRepository outbound.FriendRepository) FriendService {
	return &friendService{friendRepository}
}

func (s *friendService) CreateFriend(userID1, userID2 string) (*model.Friend, error) {
	friend, err := s.friendRepository.GetFriend(userID1, userID2)
	if err == nil {
		return nil, errors.New("friend already exists")
	}

	friend = &model.Friend{
		UserID1:  userID1,
		UserID2:  userID2,
		IsFriend: true,
		Timeline: time.Now(),
	}
	friend.ID = primitive.NewObjectID()
	friend, err = s.friendRepository.CreateFriend(friend)
	if err != nil {
		return nil, err
	}
	return friend, nil
}

func (s *friendService) GetFriend(userID1, userID2 string) (*model.Friend, error) {
	friend, err := s.friendRepository.GetFriend(userID1, userID2)
	if err != nil {
		return nil, err
	}

	return friend, nil
}

func (s *friendService) GetFriends(userID string) ([]*model.Friend, error) {
	friends, err := s.friendRepository.GetFriends(userID)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func (s *friendService) UpdateFriend(userID1, userID2 string, isFriend bool) (*model.Friend, error) {
	friend, err := s.friendRepository.GetFriend(userID1, userID2)
	if err != nil {
		return nil, err
	}

	friend.IsFriend = isFriend
	friend.Timeline = time.Now()

	friend, err = s.friendRepository.UpdateFriend(friend)
	if err != nil {
		return nil, err
	}

	return friend, nil
}

func (s *friendService) DeleteFriend(userID1, userID2 string) error {
	err := s.friendRepository.DeleteFriend(userID1, userID2)
	if err != nil {
		return err
	}

	return nil
}

func (s *friendService) GetFriendsByUserID(userID string) ([]*model.Friend, error) {
	friends, err := s.friendRepository.GetFriendsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func (s *friendService) IsFriend(userID1, userID2 string) (bool, error) {
	isFriend, err := s.friendRepository.IsFriend(userID1, userID2)
	if err != nil {
		return false, err
	}

	return isFriend, nil
}