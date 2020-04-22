package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ws(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		// TODO: handle
		return
	}

	s := newSocket(c, 0)

	s.listen()

}

// Start initiates the server
func Start(port string) {

	http.HandleFunc("/ws", ws)

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
