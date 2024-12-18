package api

import (
	"duel-masters/db"
	"duel-masters/services"
	"net/http"
)

func (api *API) getEventDeckHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return
	}

	eventUID := r.PathValue("id")

	_, err = services.ValidateEvent(eventUID)
	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": "Invalid event"})
		return
	}

	deck, err := services.GetEventDeck(user.UID, eventUID)
	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": "You are not part of this event"})
		return
	}

	d, err := services.ConvertToLegacyDeck(deck)
	if err != nil {
		write(w, 404, Json{"message": "The specified deck was not found or is not public"})
	}

	write(w, http.StatusOK, d)
}
