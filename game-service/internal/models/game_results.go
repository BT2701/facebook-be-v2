package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameResult struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID   primitive.ObjectID `bson:"session_id" json:"sessionId"`
	SymbolID    primitive.ObjectID `bson:"symbol_id" json:"symbolId"`
	Occurrences int                `bson:"occurrences" json:"occurrences"`
	OnWinline   bool               `bson:"on_winline" json:"onWinline"`
	SymbolValue float64            `bson:"symbol_value" json:"symbolValue"`
}
