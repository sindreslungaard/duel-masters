package main

import (
	"github.com/sindreslungaard/duel-masters/server"
	"github.com/sirupsen/logrus"
)

func main() {

	// TODO: change loglevel and formatter if production flag is specified
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("Starting..")
	server.Start("80")

}
