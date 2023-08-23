package migrations

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Ttl_Chatlogs_24_08_2023(db *mongo.Database) {
	logrus.Info("Creating logs collection and ttl index")

	expiresAfter := 60 * 60 * 24 * 30

	var result bson.M
	if err := db.RunCommand(context.Background(), bson.D{{"create", "logs"}}).Decode(&result); err != nil {
		logrus.Fatal(err)
	}

	index, err := db.Collection("logs").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{"ttl", 1}},
			Options: options.Index().SetExpireAfterSeconds(int32(expiresAfter)),
		},
	)

	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Successfully created ttl index with name %s", index)
}
