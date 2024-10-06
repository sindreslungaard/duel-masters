package migrations

import (
	"context"
	"duel-masters/services"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Update_Decks_23_07_2023(db *mongo.Database) {
	type Deck struct {
		UID      string   `json:"uid"`
		Owner    string   `json:"owner"`
		Name     string   `json:"name"`
		Public   bool     `json:"public"`
		Standard bool     `json:"standard"`
		Cards    []string `json:"cards"`
	}

	type UpdatedDeck struct {
		UID      string `json:"uid"`
		Owner    string `json:"owner"`
		Name     string `json:"name"`
		Public   bool   `json:"public"`
		Standard bool   `json:"standard"`
		Cards    string `json:"cards"`
	}

	decks := []Deck{}

	cur, err := db.Collection("decks").Find(context.Background(), bson.M{})

	if err != nil {
		logrus.Fatal(err)
	}

	err = cur.All(context.Background(), &decks)

	if err != nil {
		logrus.Fatal(err)
	}

	for _, deck := range decks {
		cardsMap := map[int]int{}

		for _, c := range deck.Cards {

			services.CreateIfNotExists(c)
			id, ok := services.GetCardID(c)

			if !ok {
				logrus.Fatal(fmt.Errorf("Unsuccessful migrate. Card image %s does not exist", c))
			}

			amount, ok := cardsMap[id]

			if !ok {
				cardsMap[id] = 1
			} else {
				cardsMap[id] = amount + 1
			}
		}

		digest := ""

		for id, amount := range cardsMap {
			field := strconv.Itoa(id) + "*" + strconv.Itoa(amount)
			if digest == "" {
				digest += field
			} else {
				digest += "," + field
			}
		}

		newDeck := UpdatedDeck{
			UID:      deck.UID,
			Owner:    deck.Owner,
			Name:     deck.Name,
			Public:   deck.Public,
			Standard: deck.Standard,
			Cards:    digest,
		}

		_, err := db.Collection("decks").DeleteOne(context.Background(), bson.M{"uid": newDeck.UID})

		if err != nil {
			logrus.Fatal(err)
		}

		db.Collection("decks").InsertOne(context.Background(), newDeck)
	}

}
