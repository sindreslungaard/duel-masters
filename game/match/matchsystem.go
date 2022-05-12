package match

import (
	"duel-masters/internal"
	"duel-masters/server"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
)

type MatchSystem struct {
	Matches        internal.ConcurrentDictionary[Match]
	LobbyBroadcast func(msg interface{})
}

func NewSystem(lobbyBroadcastFunc func(msg interface{})) *MatchSystem {
	return &MatchSystem{
		Matches:        internal.NewConcurrentDictionary[Match](),
		LobbyBroadcast: lobbyBroadcastFunc,
	}
}

func (s *MatchSystem) StartTicker() {

	defer internal.Recover()

	ticker := time.NewTicker(10 * time.Second) // tick every 10 seconds

	defer ticker.Stop()

	for {

		select {
		case <-ticker.C:
			{
				for _, m := range s.Matches.Iter() {
					ProcessMatch(m)
				}
			}
		}

	}

}

func ProcessMatch(m *Match) {

	defer internal.Recover()

	// Close the match if it was not started within 10 minutes of creation
	if !m.Started && m.created < time.Now().Unix()-60*10 {
		logrus.Debugf("Closing match %s", m.ID)
		m.Dispose()
	}
}

// New returns a new match object
func (s *MatchSystem) NewMatch(matchName string, hostID string, visible bool) *Match {

	id, err := shortid.Generate()

	if err != nil {
		id = uuid.New().String()
	}

	m := &Match{
		ID:                id,
		MatchName:         matchName,
		HostID:            hostID,
		spectators:        internal.NewConcurrentDictionary[Spectator](),
		persistentEffects: make(map[int]PersistentEffect),
		Turn:              1,
		Started:           false,
		Visible:           visible,

		created:     time.Now().Unix(),
		ending:      false,
		isFirstTurn: true,

		eventloop: NewEventLoop(),

		system: s,
	}

	go m.eventloop.start()

	s.Matches.Add(id, m)

	logrus.Debugf("Created match %s", id)

	return m

}

func (s *MatchSystem) UpdateMatchList() {
	s.LobbyBroadcast(MatchList(s.Matches.Iter()))
}

func MatchList(matches []*Match) server.MatchesListMessage {

	matchesMessage := make([]server.MatchMessage, 0)

	for _, match := range matches {

		if !match.Visible {
			continue
		}

		if match.ending {
			continue
		}

		if match.Player1 == nil {
			continue
		}

		if match.Player1 != nil && match.Player2 != nil && !match.Started {
			continue
		}

		matchMessage := server.MatchMessage{
			ID:      match.ID,
			P1:      match.Player1.Username,
			P1color: match.Player1.Color,
			Name:    match.MatchName,
			Started: match.Started,
		}

		if match.Player2 != nil {
			matchMessage.P2 = match.Player2.Username
			matchMessage.P2color = match.Player2.Color
		}

		matchesMessage = append(matchesMessage, matchMessage)
	}

	matchlist := server.MatchesListMessage{
		Header:  "matches",
		Matches: matchesMessage,
	}

	return matchlist

}
