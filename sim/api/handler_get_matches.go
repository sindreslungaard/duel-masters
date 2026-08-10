package api

import (
	"duel-masters/game/match"
	"net/http"
)

type matchesResponse struct {
	Matches []match.MatchSummary `json:"matches"`
}

func (api *API) getMatchesHandler(w http.ResponseWriter, r *http.Request) {
	if !authorizeServerRequest(w, r) {
		return
	}

	write(w, http.StatusOK, matchesResponse{
		Matches: api.matchSystem.MatchSummaries(),
	})
}
