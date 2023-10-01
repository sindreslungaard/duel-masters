package db

type DuelResolution string

type Duel struct {
	UID       string `json:"uid"`
	Host      string `json:"p1"`
	HostDeck  string `json:"p1deck"`
	Guest     string `json:"p2"`
	GuestDeck string `json:"p2deck"`
	Started   int64  `json:"startedAt"`
	Ended     int64  `json:"endedAt"`
	Winner    string `json:"winner"`
}

var Duels = conn().Collection("duels")
