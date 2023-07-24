package db

type Migration struct {
	Key        string
	ExecutedAt int `bson:"executed_at"`
}

var Migrations = conn().Collection("migrations")
