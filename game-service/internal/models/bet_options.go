package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BetOption struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BetSize  float64            `bson:"bet_size" json:"betSize"`
	BetLevel int                `bson:"bet_level" json:"betLevel"`
	BaseBet  float64            `bson:"base_bet" json:"baseBet"`
}
