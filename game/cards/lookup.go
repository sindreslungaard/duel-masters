package cards

import (
	"context"
	"duel-masters/db"
	"sync"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var numericLookup map[int]string
var imageLookup map[string]int
var lookupMutex sync.RWMutex
var nextId = 1

func init() {
	lookupMutex.Lock()
	defer lookupMutex.Unlock()

	cur, err := db.Cards.Find(context.TODO(), bson.M{})

	if err != nil {
		logrus.Fatal(err)
	}

	defer cur.Close(context.TODO())

	cards := make([]db.Card, 0)

	err = cur.All(context.Background(), &cards)

	if err != nil {
		logrus.Fatal(err)
	}

	for _, c := range cards {
		numericLookup[c.NumericID] = c.ImageID
		imageLookup[c.ImageID] = c.NumericID

		if c.NumericID >= nextId {
			nextId = c.NumericID + 1
		}
	}
}

func GetCardImage(id int) (string, bool) {
	lookupMutex.RLock()
	defer lookupMutex.RUnlock()
	img, ok := numericLookup[id]
	return img, ok
}

func GetCardID(image string) (int, bool) {
	lookupMutex.RLock()
	defer lookupMutex.RUnlock()
	img, ok := imageLookup[image]
	return img, ok
}

func CreateIfNotExists(image string) {
	lookupMutex.RLock()
	defer lookupMutex.RUnlock()

	_, ok := imageLookup[image]

	if ok {
		return
	}

	_, err := db.Cards.InsertOne(context.Background(), &db.Card{
		ImageID:   image,
		NumericID: nextId,
	})

	if err != nil {
		logrus.Fatal(err)
	}

	nextId = nextId + 1
}
