package api

import (
	"duel-masters/flags"
	"duel-masters/game"
	"duel-masters/game/match"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/sindreslungaard/assert"
)

type matchReqBody struct {
	HostID     string `json:"hostId" binding:"required"`
	HostDeck   string `json:"hostDeck" binding:"required"`
	GuestID    string `json:"guestId" binding:"required"`
	GuestDeck  string `json:"guestDeck" binding:"required"`
	Name       string `json:"name" binding:"max=50"`
	Visibility string `json:"visibility" binding:"required"`
	Format     string `json:"format"`
}

func (api *API) createMatchHandler(w http.ResponseWriter, r *http.Request) {
	if !flags.NewMatchesEnabled {
		write(w, http.StatusForbidden, Json{"message": "Match creation has been disabled"})
		return
	}

	var body matchReqBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	name, err := assert.Is(body.Name).MaxLen(50).String()

	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	visible := true
	if body.Visibility == "private" {
		visible = false
	}

	if name == "" {
		name = game.DefaultMatchNames[rand.Intn(len(game.DefaultMatchNames))]
	}

	format := match.FormatFromStr(body.Format)

	m := api.matchSystem.NewMatch(name, body.HostID, visible, false, format)

	write(w, http.StatusOK, m)
}
