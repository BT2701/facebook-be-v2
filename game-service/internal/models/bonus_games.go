package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BonusGame struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID primitive.ObjectID `bson:"session_id" json:"sessionId"`
	FreeSpins int                `bson:"free_spins" json:"freeSpins"`
	TotalWin  float64            `bson:"total_win" json:"totalWin"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}
