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
	// FormatID and FormatName are opaque to the simulator. It stores them on
	// the match and reports them back in match summaries and the duel result
	// webhook so the caller can identify the match's format.
	FormatID   string `json:"formatId" binding:"max=100"`
	FormatName string `json:"formatName" binding:"max=80"`
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

	if body.HostDeck == nil || body.GuestDeck == nil {
		write(w, http.StatusBadRequest, Json{"message": "hostDeck and guestDeck are required"})
		return
	}

	name, err := assert.Is(body.Name).MaxLen(50).String()

	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	formatID, err := assert.Is(strings.TrimSpace(body.FormatID)).MaxLen(100).String()

	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	formatName, err := assert.Is(strings.TrimSpace(body.FormatName)).MaxLen(80).String()

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

	format := match.FormatDescriptor{
		ID:   formatID,
		Name: formatName,
	}

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
