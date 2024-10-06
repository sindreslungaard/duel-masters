package db

type Event struct {
	UID       string   `bson:"uid"`
	Name      string   `bson:"name"`
	Organzier string   `bson:"organizer"`
	Started   int64    `bson:"startedAt"`
	Ended     int64    `bson:"endedAt"`
	Format    string   `bson:"format"`
	Sets      []string `bson:"sets"`
}

var Events = conn().Collection("events")
