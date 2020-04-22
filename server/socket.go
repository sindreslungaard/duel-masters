package server

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Socket links a ws connection to a user id and handles safe reading and writing of data
type Socket struct {
	conn   *websocket.Conn
	userID int
	mutex  *sync.Mutex
}

// newSocket creates and returns a new Socket instance
func newSocket(c *websocket.Conn, userID int) *Socket {

	s := &Socket{
		conn:   c,
		userID: userID,
		mutex:  &sync.Mutex{},
	}

	return s

}

func (s *Socket) listen() {

	defer s.Close()

	for {

		_, message, err := s.conn.ReadMessage()

		if err != nil {
			// TODO: handle
			break
		}

		log.Printf("received: %s", message)

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

}
