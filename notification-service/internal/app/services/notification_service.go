package services

import (
	"time"

	"github.com/BT2701/facebook-be-v2/shared/events"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"notification-service/internal/adapters/outbound"
	"notification-service/internal/model"
)

// NotificationService represents the notification service.
type NotificationService interface {
	CreateNotification(notification *model.Notification) (*model.Notification, error)
	GetNotification(id string) (*model.Notification, error)
	UpdateNotification(notification *model.Notification) (*model.Notification, error)
	DeleteNotification(id string) error
	GetNotificationsByUserID(userID string) ([]model.Notification, error)
	GetNotifications() ([]model.Notification, error)
	MarkAsRead(id string) error
	MarkAllAsRead(receiver string) error
	DeleteByCombo(user, receiver, post string, action int) error
	ApplyEvent(event events.NotificationEvent) error
	ApplyDeleteEvent(event events.NotificationEvent) error
}

type notificationService struct {
	notificationRepository outbound.NotificationRepository
}

// NewNotificationService creates a new notification service.
func NewNotificationService(notificationRepository outbound.NotificationRepository) NotificationService {
	return &notificationService{notificationRepository: notificationRepository}
}

func (service *notificationService) CreateNotification(notification *model.Notification) (*model.Notification, error) {
	return service.notificationRepository.CreateNotification(notification)
}

func (service *notificationService) GetNotification(id string) (*model.Notification, error) {
	return service.notificationRepository.GetNotification(id)
}

func (service *notificationService) UpdateNotification(notification *model.Notification) (*model.Notification, error) {
	return service.notificationRepository.UpdateNotification(notification)
}

func (service *notificationService) DeleteNotification(id string) error {
	return service.notificationRepository.DeleteNotification(id)
}

func (service *notificationService) GetNotificationsByUserID(userID string) ([]model.Notification, error) {
	return service.notificationRepository.GetNotificationsByUserID(userID)
}

func (service *notificationService) GetNotifications() ([]model.Notification, error) {
	return service.notificationRepository.GetNotifications()
}

func (service *notificationService) MarkAsRead(id string) error {
	return service.notificationRepository.MarkAsRead(id)
}

func (service *notificationService) MarkAllAsRead(receiver string) error {
	return service.notificationRepository.MarkAllAsRead(receiver)
}

func (service *notificationService) DeleteByCombo(user, receiver, post string, action int) error {
	return service.notificationRepository.DeleteByCombo(user, receiver, post, action)
}

func (service *notificationService) ApplyEvent(event events.NotificationEvent) error {
	post := event.Post
	if post == "" {
		post = "0"
	}
	notification := &model.Notification{
		ID:       primitive.NewObjectID(),
		User:     event.User,
		Receiver: event.Receiver,
		Post:     post,
		Content:  event.Content,
		Action_n: event.Action,
		Is_read:  0,
		Timeline: time.Now(),
	}
	_, err := service.notificationRepository.CreateNotification(notification)
	return err
}

func (service *notificationService) ApplyDeleteEvent(event events.NotificationEvent) error {
	post := event.Post
	if post == "" {
		post = "0"
	}
	return service.notificationRepository.DeleteByCombo(event.User, event.Receiver, post, event.Action)
}
