package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Story struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	UserID   string             `json:"user_id" bson:"user_id"`
	Image    string             `json:"image" bson:"image"`
	Timeline string             `json:"timeline" bson:"timeline"`
}
