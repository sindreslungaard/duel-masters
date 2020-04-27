package match

import (
	"duel-masters/server"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
)

var matches = make(map[string]*Match)
var matchesMutex = sync.Mutex{}

// Match struct
type Match struct {
	ID         string  `json:"id"`
	MatchName  string  `json:"name"`
	HostID     string  `json:"-"`
	Player1    *Player `json:"-"`
	Player2    *Player `json:"-"`
	PlayerTurn byte    `json:"-"`
}

// New returns a new match object
func New(matchName string, hostID string) *Match {

	id, err := shortid.Generate()

	if err != nil {
		id = uuid.New().String()
	}

	m := &Match{
		ID:        id,
		MatchName: matchName,
		HostID:    hostID,
	}

	matchesMutex.Lock()

	matches[id] = m

	matchesMutex.Unlock()

	logrus.Debugf("Created match %s", id)

	return m

}

// Find returns a match with the specified id, or an error
func Find(id string) (*Match, error) {

	m := matches[id]

	if m != nil {
		return m, nil
	}

	return nil, errors.New("Match does not exist")
}

// IsPlayerTurn returns a boolean based on if it is the specified player's turn
func (m *Match) IsPlayerTurn(p *Player) bool {
	if m.PlayerTurn == 1 && m.Player1 == p {
		return true
	}
	return false
}

// Parse handles websocket messages in this Hub
func (m *Match) Parse(s *server.Socket, data []byte) {

	logrus.Info(string(data))

}
