package model

import "time"

type Request struct {
	ID       string `json:"id" bson:"_id"`
	Sender  string `json:"sender" bson:"sender"`
	Receiver  string `json:"receiver" bson:"receiver"`
	Timeline time.Time `json:"timeline" bson:"timeline"`
}