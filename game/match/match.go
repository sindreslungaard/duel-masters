package match

import (
	"github.com/google/uuid"
	"github.com/ventu-io/go-shortid"
)

// Match struct
type Match struct {
	MatchName  string
	hostID     string
	Player1    *Player
	Player2    *Player
	InviteID   string
	playerTurn byte
}

// New returns a new match object
func New(matchName string, hostID string) *Match {

	id, err := shortid.Generate()

	if err != nil {
		id = uuid.New().String()
	}

	m := &Match{
		MatchName: matchName,
		hostID:    hostID,
		InviteID:  id,
	}

	return m

}

// PlayerTurn Returns the active player for this turn
func (m *Match) PlayerTurn() *Player {
	if m.playerTurn == 1 {
		return m.Player1
	}

	return m.Player2
}
