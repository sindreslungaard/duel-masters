package services

import (
	"context"
	"duel-masters/db"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var numericLookup map[int]string
var imageLookup map[string]int
var lookupMutex sync.RWMutex
var nextId = 1

func init() {
	numericLookup = map[int]string{}
	imageLookup = map[string]int{}

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
	lookupMutex.Lock()
	defer lookupMutex.Unlock()

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

	imageLookup[image] = nextId
	numericLookup[nextId] = image

	nextId = nextId + 1
}

func ConvertToLegacyDeck(deck db.Deck) (db.LegacyDeck, error) {

	cardsSlice, err := ConvertCardsStringToSlice(deck.Cards)
	if err != nil {
		return db.LegacyDeck{}, err
	}

	cardPoolSlice, err := ConvertCardsStringToSlice(deck.Cardpool)
	if err != nil {
		return db.LegacyDeck{}, err
	}

	return db.LegacyDeck{
		UID:      deck.UID,
		Owner:    deck.Owner,
		Name:     deck.Name,
		Public:   deck.Public,
		Standard: deck.Standard,
		Cards:    cardsSlice,
		Event:    deck.Event,
		Cardpool: cardPoolSlice,
	}, nil

}

func ConvertFromLegacyDeck(deck db.LegacyDeck) db.Deck {

	return db.Deck{
		UID:      deck.UID,
		Owner:    deck.Owner,
		Name:     deck.Name,
		Public:   deck.Public,
		Standard: deck.Standard,
		Cards:    ConvertCardsSliceToString(deck.Cards),
		Event:    deck.Event,
		Cardpool: ConvertCardsSliceToString(deck.Cardpool),
	}

}

func ConvertCardsSliceToString(cards []string) string {
	if cards == nil {
		return ""
	}

	cardsMap := map[int]int{}

	for _, c := range cards {
		id, ok := GetCardID(c)

		if !ok {
			logrus.Error(fmt.Errorf("Unsuccessful migrate. Card image %s does not exist", c))
			continue
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

	return digest

}

func ConvertCardsStringToSlice(cardsString string) ([]string, error) {
	if cardsString == "" {
		return []string{}, nil
	}

	cards := []string{}

	blocks := strings.Split(cardsString, ",")
	for _, block := range blocks {
		parts := strings.Split(block, "*")

		if len(parts) < 2 {
			return nil, fmt.Errorf("Invalid cards digest in deck: %s", cardsString)
		}

		cardId, err1 := strconv.Atoi(parts[0])
		amount, err2 := strconv.Atoi(parts[1])

		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("Invalid cards digest in deck: %s", cardsString)
		}

		image, ok := GetCardImage(cardId)
		if !ok {
			logrus.Error("Failed to get card image for deck")
			continue
		}

		for i := 0; i < amount; i++ {
			cards = append(cards, image)
		}

	}

	return cards, nil
}
