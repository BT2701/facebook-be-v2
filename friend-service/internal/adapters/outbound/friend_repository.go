package outbound

import (
	"context"
	"time"
	"friend-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FriendRepository interface {
	CreateFriend(friend *model.Friend) (*model.Friend, error)
	GetFriend(userID1, userID2 string) (*model.Friend, error)
	GetFriends(userID string) ([]*model.Friend, error)
	UpdateFriend(friend *model.Friend) (*model.Friend, error)
	DeleteFriend(userID1, userID2 string) error
	GetFriendsByUserID(userID string) ([]*model.Friend, error)
	IsFriend(userID1, userID2 string) (bool, error)
}

type friendRepository struct {
	collection *mongo.Collection
}

func NewFriendRepository(collection *mongo.Collection) FriendRepository {
	repo := &friendRepository{collection}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_, _ = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "userID1", Value: 1}, {Key: "userID2", Value: 1}}},
		{Keys: bson.D{{Key: "userID2", Value: 1}}},
	})
	return repo
}

func (r *friendRepository) CreateFriend(friend *model.Friend) (*model.Friend, error) {
	friend.Timeline = time.Now()
	_, err := r.collection.InsertOne(context.Background(), friend)
	if err != nil {
		return nil, err
	}
	return friend, nil
}

func (r *friendRepository) GetFriend(userID1, userID2 string) (*model.Friend, error) {
	var friend model.Friend
	err := r.collection.FindOne(context.Background(), bson.M{
		"$or": []bson.M{
			{"userID1": userID1, "userID2": userID2},
			{"userID1": userID2, "userID2": userID1},
		},
	}).Decode(&friend)
	if err != nil {
		return nil, err
	}
	return &friend, nil
}

func (r *friendRepository) GetFriends(userID string) ([]*model.Friend, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"$or": []bson.M{{"userID1": userID}, {"userID2": userID}}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var friends []*model.Friend
	for cursor.Next(context.Background()) {
		var friend model.Friend
		err := cursor.Decode(&friend)
		if err != nil {
			return nil, err
		}
		friends = append(friends, &friend)
	}
	return friends, nil
}

func (r *friendRepository) UpdateFriend(friend *model.Friend) (*model.Friend, error) {
	friend.Timeline = time.Now()
	_, err := r.collection.ReplaceOne(context.Background(), bson.M{"_id": friend.ID}, friend)
	if err != nil {
		return nil, err
	}
	return friend, nil
}

func (r *friendRepository) DeleteFriend(userID1, userID2 string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{
		"$or": []bson.M{
			{"userID1": userID1, "userID2": userID2},
			{"userID1": userID2, "userID2": userID1},
		},
	})
	return err
}

func (r *friendRepository) GetFriendsByUserID(userID string) ([]*model.Friend, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"$or": []bson.M{{"userID1": userID}, {"userID2": userID}}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var friends []*model.Friend
	for cursor.Next(context.Background()) {
		var friend model.Friend
		err := cursor.Decode(&friend)
		if err != nil {
			return nil, err
		}
		friends = append(friends, &friend)
	}
	return friends, nil
}

func (r *friendRepository) IsFriend(userID1, userID2 string) (bool, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"$or": []bson.M{{"userID1": userID1, "userID2": userID2}, {"userID1": userID2, "userID2": userID1}}})
	if err != nil {
		return false, err
	}
	defer cursor.Close(context.Background())

	return cursor.Next(context.Background()), nil
}