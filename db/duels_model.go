package db

type DuelResolution string

type Duel struct {
	UID       string `bson:"uid"`
	Host      string `bson:"p1"`
	HostDeck  string `bson:"p1deck"`
	Guest     string `bson:"p2"`
	GuestDeck string `bson:"p2deck"`
	Started   int64  `bson:"startedAt"`
	Ended     int64  `bson:"endedAt"`
	Winner    string `bson:"winner"`
}

var Duels = conn().Collection("duels")
