package outbound

import (
	"context"
	"friend-service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RequestRepository interface {
	CreateRequest(request *model.Request) (*model.Request, error)
	GetRequest(sender, receiver string) (*model.Request, error)
	GetRequests(receiver string) ([]*model.Request, error)
	UpdateRequest(request *model.Request) (*model.Request, error)
	DeleteRequest(sender, receiver string) error
	DeleteRequestByID(id string) error
	GetRequestEitherWay(sender, receiver string) (*model.Request, error)
}

type requestRepository struct {
	collection *mongo.Collection
}

func NewRequestRepository(collection *mongo.Collection) RequestRepository {
	repo := &requestRepository{collection}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	_, _ = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "sender", Value: 1}, {Key: "receiver", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "receiver", Value: 1}}},
	})
	return repo
}

func (r *requestRepository) CreateRequest(request *model.Request) (*model.Request, error) {
	request.Timeline = time.Now()
	_, err := r.collection.InsertOne(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (r *requestRepository) GetRequest(sender, receiver string) (*model.Request, error) {
	var request model.Request
	err := r.collection.FindOne(context.Background(), bson.M{"sender": sender, "receiver": receiver}).Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *requestRepository) GetRequests(receiver string) ([]*model.Request, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"receiver": receiver})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var requests []*model.Request
	for cursor.Next(context.Background()) {
		var request model.Request
		err := cursor.Decode(&request)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &request)
	}
	return requests, nil
}

func (r *requestRepository) UpdateRequest(request *model.Request) (*model.Request, error) {
	request.Timeline = time.Now()
	_, err := r.collection.ReplaceOne(context.Background(), bson.M{"_id": request.ID}, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (r *requestRepository) DeleteRequest(sender, receiver string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{
		"$or": []bson.M{
			{"sender": sender, "receiver": receiver},
			{"sender": receiver, "receiver": sender},
		},
	})
	return err
}

func (r *requestRepository) DeleteRequestByID(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

func (r *requestRepository) GetRequestEitherWay(sender, receiver string) (*model.Request, error) {
	var request model.Request
	err := r.collection.FindOne(context.Background(), bson.M{
		"$or": []bson.M{
			{"sender": sender, "receiver": receiver},
			{"sender": receiver, "receiver": sender},
		},
	}).Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}
