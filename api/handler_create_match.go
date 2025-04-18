package api

import (
	"duel-masters/db"
	"duel-masters/flags"
	"duel-masters/game"
	"duel-masters/game/match"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/sindreslungaard/assert"
)

type matchReqBody struct {
	Name       string `json:"name" binding:"max=50"`
	Visibility string `json:"visibility" binding:"required"`
	Format     string `json:"format"`
}

func (api *API) createMatchHandler(w http.ResponseWriter, r *http.Request) {
	if !flags.NewMatchesEnabled {
		write(w, http.StatusForbidden, Json{"message": "Match creation has been disabled"})
		return
	}

	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
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

	m := api.matchSystem.NewMatch(name, user.UID, visible, false, format)

	write(w, http.StatusOK, m)
}
