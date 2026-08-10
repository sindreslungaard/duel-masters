package match

import (
	"duel-masters/internal"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
)

type MatchSystem struct {
	Matches internal.ConcurrentDictionary[Match]
}

// MatchSummary is the server-to-server representation of a live match. It
// deliberately excludes decks, card state, sockets, and other private state.
type MatchSummary struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Format        string `json:"format"`
	HostID        string `json:"hostId"`
	HostUsername  string `json:"hostUsername"`
	GuestID       string `json:"guestId"`
	GuestUsername string `json:"guestUsername"`
	Visible       bool   `json:"visible"`
	Matchmaking   bool   `json:"matchmaking"`
	CreatedAt     int64  `json:"createdAt"`
}

func NewSystem() *MatchSystem {
	return &MatchSystem{
		Matches: internal.NewConcurrentDictionary[Match](),
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
func (s *MatchSystem) NewMatch(matchName string, hostID string, hostUsername string, hostDeck []string, guestID string, guestUsername string, guestDeck []string, visible bool, matchmaking bool, format Format) *Match {

	id, err := shortid.Generate()

	if err != nil {
		id = uuid.New().String()
	}

	m := &Match{
		ID:                id,
		MatchName:         matchName,
		Format:            format,
		HostID:            hostID,
		HostDeck:          hostDeck,
		GuestID:           guestID,
		GuestDeck:         guestDeck,
		spectators:        internal.NewConcurrentDictionary[Spectator](),
		persistentEffects: make(map[int]PersistentEffect),
		Turn:              1,

		PlayerSelectingToss: "",
		TossOutcome:         0,
		TossPrediction:      0,

		Started: false,
		Visible: visible,

		Matchmaking: matchmaking,
		created:     time.Now().Unix(),
		isFirstTurn: true,
		turnsPlayed: 0,

		eventloop: NewEventLoop(),

		system: s,

		hostUsername:  hostUsername,
		guestUsername: guestUsername,
	}

	go m.eventloop.start()

	s.Matches.Add(id, m)

	logrus.Debugf("Created match %s", id)

	return m

}

// MatchSummaries returns point-in-time copies of every match currently owned
// by the match system.
func (s *MatchSystem) MatchSummaries() []MatchSummary {
	matches := s.Matches.Iter()
	summaries := make([]MatchSummary, 0, len(matches))

	for _, m := range matches {
		summaries = append(summaries, m.Summary())
	}

	return summaries
}

// Summary returns the server-to-server representation of the match.
func (m *Match) Summary() MatchSummary {
	return MatchSummary{
		ID:            m.ID,
		Name:          m.MatchName,
		Format:        string(m.Format),
		HostID:        m.HostID,
		HostUsername:  m.hostUsername,
		GuestID:       m.GuestID,
		GuestUsername: m.guestUsername,
		Visible:       m.Visible,
		Matchmaking:   m.Matchmaking,
		CreatedAt:     m.created,
	}
}
