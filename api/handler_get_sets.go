package api

import (
	"net/http"
)

func (api *API) getSetsHandler(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, GetSets())
}
