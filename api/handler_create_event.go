package api

import (
	"duel-masters/db"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sindreslungaard/assert"
)

type createEventBody struct {
	Name string   `json:"name" binding:"required,min=1,max=30"`
	Sets []string `json:"sets" binding:"required"`
}

func (api *API) createEventHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return
	}

	hasPermission := false
	for _, permission := range user.Permissions {
		if permission == "admin" || permission == "chat.role.contributor" || permission == "chat.role.tournament organizer" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		write(w, http.StatusBadRequest, Json{"message": "Invalid request"})
		return
	}

	// TODO: Make sure doesn't have another ongoing tournament from same person

	var body createEventBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	name, err := assert.Is(body.Name).NotEmpty().MinLen(1).MaxLen(30).String()
	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	if len(body.Sets) > 2 || len(body.Sets) < 1 {
		write(w, http.StatusBadRequest, Json{"message": "Incorrect number of sets"})
		return
	}

	for _, set := range body.Sets {
		validset := false
		for _, s := range GetSets() {
			if set == s {
				validset = true
			}
		}
		if !validset {
			write(w, http.StatusBadRequest, Json{"message": "Invalid sets"})
			return
		}
	}

	event := db.Event{
		UID:       uuid.New().String(),
		Organzier: user.UID,
		Name:      name,
		Format:    "Sealed",
		Sets:      body.Sets,
	}

	_, err = db.Events.InsertOne(r.Context(), event)

	if err != nil {
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	write(w, http.StatusOK, Json{"message": "Event created successfully"})
	api.matchSystem.UpdateEventList()
}
