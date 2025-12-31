package migrations

import (
	"context"
	"duel-masters/db"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Migrate(conn *mongo.Database) {
	logrus.Info("Migrating database")

	type migration struct {
		key string
		fn  func(*mongo.Database)
	}

	migs := []migration{
		{key: "23_07_2023_update_decks", fn: Update_Decks_23_07_2023},
	}

	for _, m := range migs {
		count, err := db.Migrations().CountDocuments(context.Background(), bson.M{"key": m.key})

		if err != nil {
			logrus.Error(err)
		}

		if count > 0 {
			continue
		}

		m.fn(conn)

		_, err = db.Migrations().InsertOne(context.Background(), &db.Migration{
			Key:        m.key,
			ExecutedAt: int(time.Now().Unix()),
		})

		if err != nil {
			logrus.Error(err)
		}
	}
}
