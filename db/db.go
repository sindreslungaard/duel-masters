package db

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client provides a mongodb client interface
var Client *mongo.Client

// Connect connects to the database
func Connect(connectionString string) {

	clientOptions := options.Client().ApplyURI(connectionString)

	c, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		logrus.Fatal(err)
	}

	Client = c

	err = Client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	logrus.Info("Connected to MongoDB")

}
