package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reel struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Rows         int                `bson:"rows" json:"rows"`
	Columns      int                `bson:"columns" json:"columns"`
	TotalSymbols int                `bson:"total_symbols" json:"totalSymbols"`
}
