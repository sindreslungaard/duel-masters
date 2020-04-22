package server

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ws(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		// TODO: handle
		return
	}

	s := newSocket(c, 0)

	// Handle the connection in a new goroutine to free up this memory
	go s.listen()

}

// Start initiates the server
func Start(port string) {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.Dir(path.Join(dir, "public")))

	http.Handle("/", fs)
	http.HandleFunc("/ws", ws)

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
