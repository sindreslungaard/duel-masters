package main

import (
	"os"

	"duel-masters/db"
	"duel-masters/server"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {

	// TODO: change loglevel and formatter if production flag is specified
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetLevel(logrus.DebugLevel)

	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	logrus.Info("Starting..")

	db.Connect(os.Getenv("mongo_uri"), os.Getenv("mongo_name"))
	server.Start(os.Getenv("port"))

}
