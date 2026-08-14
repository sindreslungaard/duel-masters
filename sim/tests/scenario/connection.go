package scenario

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// MockConnection is a mock implementation of the Connection interface.
type MockConnection struct {
	readLimit     int64
	readDeadline  time.Time
	writeDeadline time.Time
	pongHandler   func(string) error

	MessagesWritten []string
	JSONWritten     []string
	LastJSONWritten interface{}
	Closed          bool

	mu sync.Mutex
}

// NewMockConnection creates and returns a new instance of MockConnection.
func NewMockConnection() *MockConnection {
	return &MockConnection{
		MessagesWritten: make([]string, 0),
		JSONWritten:     make([]string, 0),
	}
}

// SetReadLimit sets a dummy read limit.
func (m *MockConnection) SetReadLimit(limit int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.readLimit = limit
}

// SetReadDeadline sets a dummy read deadline.
func (m *MockConnection) SetReadDeadline(t time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.readDeadline = t
	return nil
}

// SetPongHandler sets the pong handler.
func (m *MockConnection) SetPongHandler(handler func(string) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pongHandler = handler
}

// ReadMessage returns a dummy message.
func (m *MockConnection) ReadMessage() (messageType int, p []byte, err error) {
	return 1, []byte("mock message"), nil
}

// SetWriteDeadline sets a dummy write deadline.
func (m *MockConnection) SetWriteDeadline(t time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.writeDeadline = t
	return nil
}

// WriteMessage simulates writing a message.
func (m *MockConnection) WriteMessage(messageType int, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MessagesWritten = append(m.MessagesWritten, string(data))
	return nil
}

// WriteJSON simulates writing a JSON message.
func (m *MockConnection) WriteJSON(v interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	m.LastJSONWritten = v
	m.JSONWritten = append(m.JSONWritten, string(data))
	return err
}

func (m *MockConnection) JSONWriteCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.JSONWritten)
}

func (m *MockConnection) JSONMessagesSince(start int) []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if start >= len(m.JSONWritten) {
		return []string{}
	}

	result := make([]string, len(m.JSONWritten[start:]))
	copy(result, m.JSONWritten[start:])
	return result
}

// Close simulates closing the connection.
func (m *MockConnection) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Closed {
		return errors.New("connection already closed")
	}
	m.Closed = true
	return nil
}
