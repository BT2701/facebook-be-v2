package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID   string             `bson:"user_id" json:"user_id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password,omitempty" json:"-"`
	Balance  float64            `bson:"balance" json:"balance"`
}
