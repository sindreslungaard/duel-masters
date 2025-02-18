package api

import (
	"context"
	"duel-masters/db"
	"duel-masters/services"
	"net/http"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	Status_Joined   = "joined"
	Status_Playable = "playable"
)

type EventForPlayer struct {
	db.Event
	Status string `json:"status"`
}
type EventsResponse struct {
	EventsForPlayer []EventForPlayer `json:"events"`
}

func (api *API) getEventsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return
	}

	cur, err := db.Events.Find(r.Context(), bson.M{
		"format":  "Sealed",
		"endedAt": 0,
	})

	if err != nil {
		logrus.Error(err)
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	var events []db.Event
	var eventIds []string
	for cur.Next(context.TODO()) {
		var event db.Event
		if err := cur.Decode(&event); err != nil {
			continue
		}
		events = append(events, event)
		eventIds = append(eventIds, event.UID)
	}

	var response EventsResponse
	if len(eventIds) == 0 {
		write(w, http.StatusOK, response)
		return
	}

	cur, err = db.Decks.Find(
		r.Context(),
		bson.M{"owner": user.UID, "event": bson.M{"$in": eventIds}},
	)

	if err != nil {
		logrus.Error(err)
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	eventDeck := make(map[string]db.Deck)
	for cur.Next(context.TODO()) {
		var deck db.Deck
		if err := cur.Decode(&deck); err != nil {
			continue
		}
		eventDeck[deck.Event] = deck
	}

	for _, event := range events {
		var efp EventForPlayer
		efp.Event = event
		deck, ok := eventDeck[event.UID]
		if ok {
			cards, _ := services.ConvertCardsStringToSlice(deck.Cards)
			if len(cards) == 40 {
				efp.Status = Status_Playable
			} else {
				efp.Status = Status_Joined
			}
		}
		response.EventsForPlayer = append(response.EventsForPlayer, efp)
	}

	write(w, http.StatusOK, response)
}
