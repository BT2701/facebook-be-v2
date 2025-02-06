package models


type Reel struct {
	Rows         int                `bson:"rows" json:"rows"`
	Columns      int                `bson:"columns" json:"columns"`
}
