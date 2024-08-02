package api

import (
	"duel-masters/db"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (api *API) deleteDeckHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusForbidden, Json{"message": "Match creation has been disabled"})
		return
	}

	deckUID := r.PathValue("id")

	result, err := db.Decks.DeleteOne(
		r.Context(),
		bson.M{"uid": deckUID, "owner": user.UID},
	)

	if err != nil {
		write(w, http.StatusNotFound, Json{"message": "Deck not found"})
		return
	}

	if result.DeletedCount < 1 {
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	write(w, http.StatusOK, Json{"message": "Successfully deleted the deck"})
}
