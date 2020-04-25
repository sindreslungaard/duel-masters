package server

import (
	"log"
	"sync"

	"duel-masters/db"

	"github.com/gorilla/websocket"
)

var sockets = make(map[*Socket]int)
var socketsMutex = sync.Mutex{}

// Socket links a ws connection to a user id and handles safe reading and writing of data
type Socket struct {
	conn  *websocket.Conn
	user  db.User
	mutex *sync.Mutex
}

// newSocket creates and returns a new Socket instance
func newSocket(c *websocket.Conn, user db.User) *Socket {

	s := &Socket{
		conn:  c,
		user:  user,
		mutex: &sync.Mutex{},
	}

	return s

}

func (s *Socket) listen() {

	defer func() {

		socketsMutex.Lock()
		//delete(sockets, s.user.UID)
		socketsMutex.Unlock()

		s.Close()

	}()

	for {

		_, message, err := s.conn.ReadMessage()

		if err != nil {
			// TODO: handle
			break
		}

		Parse(s, message)

	}

}

// Send sends a struct v to the client
func (s *Socket) Send(v interface{}) {

	s.mutex.Lock()
	s.conn.WriteJSON(v)
	s.mutex.Unlock()

}

// Close closes the client connection
func (s *Socket) Close() {

	s.conn.Close()
	log.Println("Closed a connection")

}
