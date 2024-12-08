package api

import (
	"context"
	"duel-masters/db"
	"duel-masters/services"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (api *API) getDeckHandler(w http.ResponseWriter, r *http.Request) {
	deckUID := r.PathValue("id")

	var deck db.Deck

	err := db.Decks.FindOne(
		r.Context(),
		bson.M{"uid": deckUID, "public": true},
	).Decode(&deck)

	if err != nil {
		write(w, 404, Json{"message": "The specified deck was not found or is not public"})
		return
	}

	var user db.User

	err = db.Users.FindOne(
		context.Background(),
		bson.M{"uid": deck.Owner},
	).Decode(&user)

	if err != nil {
		write(w, 404, Json{"message": "The specified deck was not found or is not public"})
		return
	}

	deck.Owner = user.Username

	d, err := services.ConvertToLegacyDeck(deck)

	if err != nil {
		write(w, 404, Json{"message": "The specified deck was not found or is not public"})
	}

	write(w, http.StatusOK, d)
}
