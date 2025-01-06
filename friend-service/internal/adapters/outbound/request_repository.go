package outbound

import (
	"context"
	"friend-service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RequestRepository interface {
	CreateRequest(request *model.Request) (*model.Request, error)
	GetRequest(sender, receiver string) (*model.Request, error)
	GetRequests(receiver string) ([]*model.Request, error)
	UpdateRequest(request *model.Request) (*model.Request, error)
	DeleteRequest(sender, receiver string) error
}

type requestRepository struct {
	collection *mongo.Collection
}

func NewRequestRepository(collection *mongo.Collection) RequestRepository {
	return &requestRepository{collection}
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
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"sender": sender, "receiver": receiver})
	if err != nil {
		return err
	}
	return nil
}
