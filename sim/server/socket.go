package server

import (
	"sync"
	"sync/atomic"
	"time"

	"duel-masters/internal"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Connection interface {
	SetReadLimit(int64)
	SetReadDeadline(t time.Time) error
	SetPongHandler(func(string) error)
	ReadMessage() (messageType int, p []byte, err error)
	SetWriteDeadline(t time.Time) error
	WriteMessage(messageType int, data []byte) error
	WriteJSON(v interface{}) error
	Close() error
}

// Socket links a ws connection to a user id and handles safe reading and writing of data
type Socket struct {
	conn  Connection
	User  User
	hub   Hub
	ready bool
	mutex *sync.Mutex
	// closed and lost are written when the connection goes away, which happens
	// on the socket's own goroutine while other goroutines are still sending on
	// it, so they must be race free.
	closed atomic.Bool
	lost   atomic.Bool
}

// NewSocket creates and returns a new Socket instance
func NewSocket(c Connection, hub Hub, userID string, username string) *Socket {
	var user User
	user.UID = userID
	user.Username = username

	s := &Socket{
		conn:  c,
		hub:   hub,
		ready: true,
		mutex: &sync.Mutex{},
		User:  user,
	}

	logrus.Debugf("Opened a connection")

	return s

}

// Ready returns true or false based on if the socket is ready or not
func (s *Socket) Ready() bool {
	return s.ready
}

// Listen sets up reader and writer for the socket
func (s *Socket) Listen() {

	s.conn.SetReadLimit(maxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(pongWait))
	s.conn.SetPongHandler(func(string) error { s.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	defer s.Close()

	go s.handlePing()

	for {

		_, message, err := s.conn.ReadMessage()

		if err != nil {
			return
		}

		s.hub.Parse(s, message)

	}

}

func (s *Socket) handlePing() {

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("recovered from handlePing: %v", r)
		}
	}()

	for {

		if s.isGone() {
			return
		}

		select {
		case <-ticker.C:
			s.mutex.Lock()
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := s.conn.WriteMessage(websocket.PingMessage, nil)
			s.mutex.Unlock()

			if err != nil {
				if !s.isGone() {
					s.conn.Close()
				}
				return
			}
		}
	}

}

// Send sends a struct v to the client
func (s *Socket) Send(v any) {

	if s.isGone() {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from panic in socket Send. %v", r)
			return
		}
	}()

	s.mutex.Lock()
	s.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := s.conn.WriteJSON(v); err != nil {
		logrus.Debug(err)
	}
	s.mutex.Unlock()

}

// Close closes the client connection
func (s *Socket) Close() {

	defer internal.Recover()

	if s.closed.Swap(true) {
		return
	}

	s.hub.OnSocketClose(s)

	if s.conn != nil {
		s.conn.Close()
	}

	logrus.Debug("Closed a connection")

}

func (s *Socket) IsClosed() bool {
	return s.closed.Load()
}

// isGone reports whether the connection can no longer carry messages.
func (s *Socket) isGone() bool {
	return s.closed.Load() || s.lost.Load()
}

func (s *Socket) Warn(msg string) {
	s.Send(WarningMessage{
		Header:  "warn",
		Message: msg,
	})
}
