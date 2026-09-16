package outbound

import (
	"context"
	"log"
	"notification-service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationRepository interface {
	CreateNotification(notification *model.Notification) (*model.Notification, error)
	GetNotification(id string) (*model.Notification, error)
	UpdateNotification(notification *model.Notification) (*model.Notification, error)
	DeleteNotification(id string) error
	GetNotificationsByUserID(userID string) ([]model.Notification, error)
	GetNotifications() ([]model.Notification, error)
	MarkAsRead(id string) error
	MarkAllAsRead(receiver string) error
	DeleteByCombo(user, receiver, post string, action int) error
}

type notificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository(collection *mongo.Collection) NotificationRepository {
	repo := &notificationRepository{collection: collection}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_, _ = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "receiver", Value: 1}, {Key: "timeline", Value: -1}}},
		{Keys: bson.D{{Key: "user", Value: 1}, {Key: "receiver", Value: 1}, {Key: "post", Value: 1}, {Key: "action_n", Value: 1}}},
	})
	return repo
}

func (repo *notificationRepository) CreateNotification(notification *model.Notification) (*model.Notification, error) {
	if notification.Action_n == 1 || notification.Action_n == 2 || notification.Action_n == 3 || notification.Action_n == 4 {
		filter := bson.M{
			"user":     notification.User,
			"receiver": notification.Receiver,
			"post":     notification.Post,
			"action_n": notification.Action_n,
		}
		update := bson.M{"$setOnInsert": notification}
		opts := options.Update().SetUpsert(true)
		result, err := repo.collection.UpdateOne(context.Background(), filter, update, opts)
		if err != nil {
			return nil, err
		}
		if result.UpsertedID == nil {
			var existing model.Notification
			if err := repo.collection.FindOne(context.Background(), filter).Decode(&existing); err == nil {
				return &existing, nil
			}
		}
		return notification, nil
	}

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
	opts := options.Find().SetSort(bson.D{{Key: "timeline", Value: -1}}).SetLimit(50)
	cursor, err := repo.collection.Find(context.Background(), bson.M{"receiver": userID}, opts)
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

func (repo *notificationRepository) MarkAsRead(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = repo.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": bson.M{"is_read": 1}})
	return err
}

func (repo *notificationRepository) MarkAllAsRead(receiver string) error {
	_, err := repo.collection.UpdateMany(context.Background(), bson.M{"receiver": receiver}, bson.M{"$set": bson.M{"is_read": 1}})
	return err
}

func (repo *notificationRepository) DeleteByCombo(user, receiver, post string, action int) error {
	_, err := repo.collection.DeleteMany(context.Background(), bson.M{
		"user":     user,
		"receiver": receiver,
		"post":     post,
		"action_n": action,
	})
	return err
}

