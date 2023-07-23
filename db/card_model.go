package db

type Card struct {
	NumericID int    `json:"numeric_id"`
	ImageID   string `json:"image_id"`
}

var Cards = conn().Collection("cards")
