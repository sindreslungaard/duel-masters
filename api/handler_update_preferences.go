package api

import (
	"duel-masters/db"
	"encoding/json"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type preferencesReqBody struct {
	Playmat string `json:"playmat"`
}

func (api *API) updatePreferencesHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return
	}

	var body preferencesReqBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	if body.Playmat != "" && !strings.HasPrefix(body.Playmat, "https://i.imgur.com/") {
		write(w, http.StatusBadRequest, Json{"message": "Playmat images must be uploaded to imgur and the url must start with https://i.imgur.com/. Make sure the link includes the file extension (.png, .jpg)"})
		return
	}

	db.Users.UpdateOne(r.Context(), bson.M{
		"uid": user.UID,
	}, bson.M{"$set": bson.M{
		"playmat": body.Playmat,
	}})

	write(w, http.StatusOK, Json{"message": "Successfully saved your preferences"})
}
