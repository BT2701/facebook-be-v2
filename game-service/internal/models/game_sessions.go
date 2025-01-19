package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type GameSession struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PlayerID       primitive.ObjectID `bson:"player_id" json:"playerId"`
	BetOptionID    primitive.ObjectID `bson:"bet_option_id" json:"betOptionId"`
	TotalBetAmount float64            `bson:"total_bet_amount" json:"totalBetAmount"`
	TotalWin       float64            `bson:"total_win" json:"totalWin"`
	Status         string             `bson:"status" json:"status"` // active or finished
	CreatedAt      time.Time          `bson:"created_at" json:"createdAt"`
}
