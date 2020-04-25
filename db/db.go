package db

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var conn *mongo.Database

// Connect connects to the database
func Connect(connectionString string, dbName string) {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))

	err = client.Ping(context.TODO(), readpref.Primary())

	if err != nil {
		logrus.Fatal(err)
	}

	conn = client.Database(dbName)

	logrus.Info("Connected to database")

}

// Collection returns a mongodb collection handle
func Collection(collectionName string) *mongo.Collection {
	return conn.Collection(collectionName)
}
