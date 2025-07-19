package db

import "go.mongodb.org/mongo-driver/mongo"

type BanType string

const (
	UserBan BanType = "user"
	IPBan   BanType = "ip"
)

type Ban struct {
	Type  BanType `json:"type"`
	Value string  `json:"value"`
}

func Bans() *mongo.Collection {
	return conn().Collection("bans")
}
