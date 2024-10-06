package api

import (
	"duel-masters/db"
	"duel-masters/services"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (api *API) joinEventHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = services.GetEventDeck(user.UID, eventUID)
	if err == nil {
		write(w, http.StatusBadRequest, Json{"message": "You already joined this event"})
		return
	}

	// Generate card pool
	set := event.Sets[0]
	var cardPool []string
	for i := 0; i < 8; i++ {
		if i == 4 && len(event.Sets) == 2 {
			set = event.Sets[1]
		}
		pack := GeneratePack(set)
		for _, cards := range pack {
			cardPool = append(cardPool, cards...)
		}
	}

	newDeck := db.LegacyDeck{
		UID:      uuid.New().String(),
		Owner:    user.UID,
		Name:     fmt.Sprintf("%s_%s", event.Name, user.Username),
		Public:   false,
		Standard: false,
		Cards:    nil,
		Event:    eventUID,
		Cardpool: cardPool,
	}

	_, err = db.Decks.InsertOne(r.Context(), services.ConvertFromLegacyDeck(newDeck))
	if err != nil {
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	write(w, http.StatusOK, Json{"message": "You joined the event"})
}
