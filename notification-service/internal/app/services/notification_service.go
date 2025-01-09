package services

import (
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
