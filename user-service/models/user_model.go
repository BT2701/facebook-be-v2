package models

import "time"

type User struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	Name         string    `bson:"name" json:"name"`
	Email        string    `bson:"email" json:"email"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	Password     string    `bson:"password" json:"password"`
	Description  string    `bson:"description" json:"description"`
	Avatar       string    `bson:"avatar" json:"avatar"`
	Birthday     time.Time `bson:"birthday" json:"birthday"`
	Gender       string    `bson:"gender" json:"gender"`
	Phone        string    `bson:"phone" json:"phone"`
	Is_online    int       `bson:"is_online" json:"is_online"`
	Last_active  time.Time `bson:"last_active" json:"last_active"`
	Address      string    `bson:"address" json:"address"`
	Social       string    `bson:"social" json:"social"`
	Education    string    `bson:"education" json:"education"`
	Relationship string    `bson:"relationship" json:"relationship"`
}
