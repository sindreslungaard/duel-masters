package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sindreslungaard/duel-masters/db"
	"github.com/sindreslungaard/duel-masters/server"
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

	db.Connect(os.Getenv("db_host"), os.Getenv("db_port"), os.Getenv("db_user"), os.Getenv("db_pass"), os.Getenv("db_name"))
	server.Start(os.Getenv("port"))

}
