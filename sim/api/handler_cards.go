package api

import (
	"net/http"
)

func (api *API) getCardsHandler(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, GetCache())
}
