package api

import (
	"duel-masters/flags"
	"duel-masters/game/match"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"

	"github.com/sindreslungaard/assert"
)

type matchReqBody struct {
	HostID        string   `json:"hostId" binding:"required"`
	HostUsername  string   `json:"hostUsername" binding:"required"`
	HostDeck      []string `json:"hostDeck" binding:"required"`
	GuestID       string   `json:"guestId" binding:"required"`
	GuestUsername string   `json:"guestUsername" binding:"required"`
	GuestDeck     []string `json:"guestDeck" binding:"required"`
	Name          string   `json:"name" binding:"max=50"`
	Visibility    string   `json:"visibility" binding:"required"`
	Format        string   `json:"format"`
}

var defaultMatchNames = []string{
	"Kettou Da!",
	"I challenge you!",
	"Ikuzo!",
	"I'm ready!",
	"Koi!",
	"Bring it on!",
}

func (api *API) createMatchHandler(w http.ResponseWriter, r *http.Request) {
	if !authorizeServerRequest(w, r) {
		return
	}

	if !flags.NewMatchesEnabled {
		write(w, http.StatusForbidden, Json{"message": "Match creation has been disabled"})
		return
	}

	var body matchReqBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	body.HostUsername = strings.TrimSpace(body.HostUsername)
	body.GuestUsername = strings.TrimSpace(body.GuestUsername)
	if body.HostUsername == "" || body.GuestUsername == "" {
		write(w, http.StatusBadRequest, Json{"message": "hostUsername and guestUsername are required"})
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
		name = defaultMatchNames[rand.Intn(len(defaultMatchNames))]
	}

	format := match.FormatFromStr(body.Format)

	m := api.matchSystem.NewMatch(
		name,
		body.HostID,
		body.HostUsername,
		body.HostDeck,
		body.GuestID,
		body.GuestUsername,
		body.GuestDeck,
		visible,
		false,
		format,
	)

	write(w, http.StatusOK, m)
}
