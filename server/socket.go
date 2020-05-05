package server

import (
	"log"
	"sync"

	"duel-masters/db"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var sockets = make(map[*Socket]Hub)
var socketsMutex = sync.Mutex{}

// Socket links a ws connection to a user id and handles safe reading and writing of data
type Socket struct {
	conn  *websocket.Conn
	User  db.User
	hub   Hub
	ready bool
	mutex *sync.Mutex
}

// NewSocket creates and returns a new Socket instance
func NewSocket(c *websocket.Conn, hub Hub) *Socket {

	s := &Socket{
		conn:  c,
		hub:   hub,
		ready: false,
		mutex: &sync.Mutex{},
	}

	socketsMutex.Lock()
	sockets[s] = hub
	socketsMutex.Unlock()

	return s

}

// Listen sets up reader and writer for the socket
func (s *Socket) Listen() {

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

		if !s.ready {

			// Look for authorization token as the first message
			u, err := db.GetUserForToken(string(message))

			if err != nil {
				continue
			}

			s.User = u
			s.ready = true

			s.Send(Message{Header: "hello"})

			continue

		}

		s.hub.Parse(s, message)

	}

}

// Send sends a struct v to the client
func (s *Socket) Send(v interface{}) {

	s.mutex.Lock()
	if err := s.conn.WriteJSON(v); err != nil {
		logrus.Debug(err)
	}
	s.mutex.Unlock()

}

// Close closes the client connection
func (s *Socket) Close() {

	s.conn.Close()
	log.Println("Closed a connection")

}
