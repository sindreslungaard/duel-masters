package api

import (
	"context"
	"duel-masters/db"
	"duel-masters/services"
	"net/http"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func (api *API) getDecksHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return
	}

	cur, err := db.Decks.Find(r.Context(), bson.M{
		"owner": user.UID,
		"event": nil,
	})

	if err != nil {
		logrus.Error(err)
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	defer cur.Close(context.TODO())

	decks := make([]db.Deck, 0)

	for cur.Next(context.TODO()) {

		var deck db.Deck

		if err := cur.Decode(&deck); err != nil {
			continue
		}

		decks = append(decks, deck)

	}

	legacyDecks := []db.LegacyDeck{}

	for _, deck := range decks {
		legacyDeck, err := services.ConvertToLegacyDeck(deck)
		if err != nil {
			continue
		}
		legacyDecks = append(legacyDecks, legacyDeck)
	}

	write(w, http.StatusOK, legacyDecks)
}
