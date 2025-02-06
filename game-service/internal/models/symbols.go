package models

type (
	Symbol struct {
		ID     int                `json:"id" bson:"id"`
		Name   string             `json:"name" bson:"name"`
		Values map[string]float64 `json:"values" bson:"values"`
		Color  string             `json:"color" bson:"color"`
	}
)

type GameSymbolData struct {
	Symbols []Symbol `json:"symbols" bson:"symbols"`
}
