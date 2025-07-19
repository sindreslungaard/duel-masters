package db

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const ConnectionStringEnv = "mongo_uri"
const DatabaseNameEnv = "mongo_name"

var connection *mongo.Database

func connect() {
	connectionString := os.Getenv(ConnectionStringEnv)
	dbName := os.Getenv(DatabaseNameEnv)

	if connectionString == "" || dbName == "" {
		logrus.Fatal(fmt.Sprintf("Missing %s or %s environment variables", ConnectionStringEnv, DatabaseNameEnv))
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))

	if err != nil {
		logrus.Fatal(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())

	if err != nil {
		logrus.Fatal(err)
	}

	connection = client.Database(dbName)

	logrus.Info("Connected to database")
}

func conn() *mongo.Database {
	if connection == nil {
		connect()
	}
	return connection
}

// GetUserForToken returns a user from the authorization header or returns an error
func GetUserForToken(token string) (User, error) {

	var user User

	if err := Users().FindOne(context.TODO(), bson.M{"sessions": bson.M{"$elemMatch": bson.M{"token": token}}}).Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil

}

func Connection() *mongo.Database {
	return connection
}
