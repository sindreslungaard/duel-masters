package main

import (
	"math/rand"
	"os"
	"time"

	"duel-masters/api"
	"duel-masters/db"
	"duel-masters/game/cards"
	"duel-masters/game/match"

	"github.com/sirupsen/logrus"
)

func main() {

	// TODO: change loglevel and formatter if production flag is specified
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetLevel(logrus.DebugLevel)

	rand.Seed(time.Now().UnixNano())

	logrus.Info("Starting..")

	for uid, ctor := range cards.DM01 {
		match.AddCard(uid, ctor)
	}

	api.CreateCardCache()

	db.Connect(os.Getenv("mongo_uri"), os.Getenv("mongo_name"))

	api.Start(os.Getenv("port"))

}
