package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Symbol struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"` // scatter, wild, p1, etc.
	BaseValue float64            `bson:"base_value" json:"baseValue"`
	IsScatter bool               `bson:"is_scatter" json:"isScatter"`
	IsWild    bool               `bson:"is_wild" json:"isWild"`
}
