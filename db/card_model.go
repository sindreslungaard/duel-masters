package db

import "go.mongodb.org/mongo-driver/mongo"

type Card struct {
	NumericID int    `json:"numeric_id"`
	ImageID   string `json:"image_id"`
}

func Cards() *mongo.Collection {
	return conn().Collection("cards")
}
