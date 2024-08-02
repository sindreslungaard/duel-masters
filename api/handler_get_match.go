package api

import (
	"net/http"
)

func (api *API) getMatchHandler(w http.ResponseWriter, r *http.Request) {
	m, ok := api.matchSystem.Matches.Find(r.PathValue("id"))

	if !ok {
		write(w, http.StatusNotFound, Json{"message": "Match not found"})
		return
	}

	write(w, http.StatusOK, Json{"name": m.MatchName, "host": m.HostID, "started": m.Started})
}
