package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Balance  float64            `bson:"balance" json:"balance"`
}
