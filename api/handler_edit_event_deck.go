package api

import (
	"duel-masters/db"
	"duel-masters/services"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type createEventDeckBody struct {
	Cards []string `json:"cards" binding:"required"`
}

func (api *API) editEventDeckHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return
	}

	eventUID := r.PathValue("id")
	event, err := services.ValidateEvent(eventUID)
	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": "Invalid event"})
		return
	}

	existingDeck, err := services.GetEventDeck(user.UID, event.UID)
	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": "You are not part of this event."})
		return
	}

	var body createEventDeckBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	// TODO: Should be based on the event rules
	if body.Cards == nil || len(body.Cards) != 40 {
		write(w, http.StatusBadRequest, Json{"message": "Deck must consist of 40 cards"})
		return
	}

	cardPool, err := services.ConvertCardsStringToSlice(existingDeck.Cardpool)
	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": "Internal error, contact admin"})
		return
	}

	for _, card := range body.Cards {
		cardFound := false
		for i, cardFromPool := range cardPool {
			if cardFromPool == card {
				cardPool = append(cardPool[:i], cardPool[i+1:]...)
				cardFound = true
				break
			}
		}
		if !cardFound {
			write(w, http.StatusBadRequest, Json{"message": "You are using cards that are not in your list"})
			return
		}
	}

	_, err = db.Decks.UpdateOne(
		r.Context(),
		bson.M{"uid": existingDeck.UID, "owner": user.UID},
		bson.M{"$set": bson.M{"cards": services.ConvertCardsSliceToString(body.Cards)}},
	)
	if err != nil {
		logrus.Error(err)
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	write(w, http.StatusOK, Json{"message": "Deck saved successfully"})
}
