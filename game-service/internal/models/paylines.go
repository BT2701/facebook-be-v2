package models

type (
	Payline struct {
		ID string `json:"id" bson:"id"`
		Pattern [][]int `json:"pattern,omitempty" bson:"pattern,omitempty"`
		Adjustable bool `json:"adjustable" bson:"adjustable"`
		Direction string `json:"direction,omitempty" bson:"direction,omitempty"`
		Index int `json:"index" bson:"index"`
	}
)

type GamePaylines struct {
	Paylines []Payline `json:"paylines" bson:"paylines"`
}