package api

import (
	"duel-masters/db"
	"duel-masters/services"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func (api *API) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return
	}

	eventUID := r.PathValue("id")
	event, err := services.ValidateEvent(eventUID)
	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	if event.Organzier != user.UID {
		write(w, http.StatusBadRequest, Json{"message": "Not Allowed"})
		return
	}

	_, err = db.Events.UpdateOne(
		r.Context(),
		bson.M{"uid": event.UID},
		bson.M{"$set": bson.M{"endedAt": time.Now().Unix()}},
	)
	if err != nil {
		logrus.Error(err)
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	write(w, http.StatusOK, Json{"message": "Event terminated"})
	api.matchSystem.UpdateEventList()
}
