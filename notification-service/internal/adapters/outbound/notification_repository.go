package outbound

import (
	"context"
	"log"
	"notification-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository interface {
	CreateNotification(notification *model.Notification) (*model.Notification, error)
	GetNotification(id string) (*model.Notification, error)
	UpdateNotification(notification *model.Notification) (*model.Notification, error)
	DeleteNotification(id string) error
	GetNotificationsByUserID(userID string) ([]model.Notification, error)
	GetNotifications() ([]model.Notification, error)
}

type notificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository(collection *mongo.Collection) NotificationRepository {
	return &notificationRepository{collection: collection}
}

func (repo *notificationRepository) CreateNotification(notification *model.Notification) (*model.Notification, error) {
	_, err := repo.collection.InsertOne(context.Background(), notification)
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (repo *notificationRepository) GetNotification(id string) (*model.Notification, error) {
	var notification model.Notification
	err := repo.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&notification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &notification, nil
}

func (repo *notificationRepository) UpdateNotification(notification *model.Notification) (*model.Notification, error) {
	_, err := repo.collection.ReplaceOne(context.Background(), bson.M{"_id": notification.ID}, notification)
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (repo *notificationRepository) DeleteNotification(id string) error {
	_, err := repo.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (repo *notificationRepository) GetNotificationsByUserID(userID string) ([]model.Notification, error) {
	var notifications []model.Notification
	cursor, err := repo.collection.Find(context.Background(), bson.M{"receiver": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var notification model.Notification
		if err := cursor.Decode(&notification); err != nil {
			log.Println("Error decoding notification:", err)
			continue
		}
		notifications = append(notifications, notification)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (repo *notificationRepository) GetNotifications() ([]model.Notification, error) {
	var notifications []model.Notification
	cursor, err := repo.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var notification model.Notification
		if err := cursor.Decode(&notification); err != nil {
			log.Println("Error decoding notification:", err)
			continue
		}
		notifications = append(notifications, notification)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

