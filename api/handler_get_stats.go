package api

import (
	"duel-masters/db"
	"net/http"
)

func (api *API) getStatsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return
	}

	write(w, http.StatusOK, Json{
		"games_lost":         user.GamesLost,
		"games_won":          user.GamesWon,
		"total_games_played": user.TotalGamesPlayed,
		"win_rate":           int(user.WinRate),
	})
}
