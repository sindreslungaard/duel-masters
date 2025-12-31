package db

import "go.mongodb.org/mongo-driver/mongo"

type Migration struct {
	Key        string
	ExecutedAt int `bson:"executed_at"`
}

func Migrations() *mongo.Collection {
	return conn().Collection("migrations")
}
