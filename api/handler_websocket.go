package api

import (
	"duel-masters/server"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (api *API) websocketHandler(w http.ResponseWriter, r *http.Request) {
	hubID := r.Header.Get("hub")

	var hub server.Hub

	if hubID == "lobby" {
		hub = api.lobby
	} else {
		m, ok := api.matchSystem.Matches.Find(hubID)

		if !ok {
			write(w, http.StatusNotFound, Json{"message": "Match not found"})
			return
		}

		hub = m
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	s := server.NewSocket(conn, hub)

	// Handle the connection in a new goroutine to free up this memory
	go s.Listen()
}
