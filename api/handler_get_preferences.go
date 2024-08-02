package api

import (
	"duel-masters/db"
	"net/http"
)

func (api *API) getPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return
	}

	write(w, http.StatusOK, Json{
		"playmat": user.Playmat,
	})
}
