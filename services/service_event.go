package services

import (
	"context"
	"duel-masters/db"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

// Checks if an event exists and it's still allowed to play in it
func ValidateEvent(id string) (db.Event, error) {
	var event db.Event
	err := db.Events.FindOne(
		context.TODO(),
		bson.M{"uid": id, "endedAt": 0},
	).Decode(&event)
	return event, err
}

func GetEventDeck(userUID string, eventUID string) (db.Deck, error) {
	var deck db.Deck
	err := db.Decks.FindOne(
		context.TODO(),
		bson.M{"owner": userUID, "event": eventUID},
	).Decode(&deck)

	return deck, err
}

func CanPlayerPlayEvent(userUID string, eventUID string) (bool, error) {
	_, err := ValidateEvent(eventUID)
	if err != nil {
		return false, err
	}

	deck, err := GetEventDeck(userUID, eventUID)
	if err != nil {
		return false, err
	}

	legacyDeck, err := ConvertToLegacyDeck(deck)
	if err != nil {
		return false, errors.New("internal error")
	}

	// Check event req in the future
	if len(legacyDeck.Cards) != 40 {
		return false, errors.New("invalid deck")
	}

	return true, nil
}
