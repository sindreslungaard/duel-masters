package game

import (
	"duel-masters/game/match"
	"duel-masters/internal"
	"duel-masters/server"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

var Matchmaking = &Matchmaker{
	requests:    internal.NewConcurrentDictionary[MatchRequest](),
	broadcaster: func(msg interface{}) { logrus.Warn("Use of default matchmaking broadcaster") },
}

type MatchUser struct {
	ID       string
	Username string
	SocketID string
}

type MatchRequest struct {
	ID           string
	Host         MatchUser
	Guest        *internal.Option[MatchUser]
	Name         string
	Format       match.Format
	BlockedUsers map[string]bool
}

type Matchmaker struct {
	sync.RWMutex
	requests    internal.ConcurrentDictionary[MatchRequest]
	broadcaster func(msg interface{})
}

func (m *Matchmaker) SetBroadcaster(f func(msg interface{})) {
	m.broadcaster = f
}

func (m *Matchmaker) NewRequest(s *server.Socket, name string, format match.Format) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.requests.Find(s.User.UID); ok {
		return fmt.Errorf("you already have a duel request open")
	}

	r := &MatchRequest{
		ID: s.User.UID,
		Host: MatchUser{
			ID:       s.User.UID,
			Username: s.User.Username,
			SocketID: s.UID,
		},
		Guest:        internal.NewOption[MatchUser](),
		Name:         name,
		Format:       format,
		BlockedUsers: map[string]bool{},
	}

	m.requests.Add(s.User.UID, r)

	return nil
}

func (m *Matchmaker) Remove(id string) {
	m.Lock()
	defer m.Unlock()

	m.requests.Remove(id)
}

func (m *Matchmaker) SocketClosed(id string) {
	for _, r := range m.requests.Iter() {
		if r.Host.SocketID == id {
			m.Remove(r.ID)
		}
	}
}

func (m *Matchmaker) Join(s *server.Socket, id string) error {
	// check if guest in another request
	for _, r := range m.requests.Iter() {
		if guest, ok := r.Guest.Unwrap(); ok && guest.ID == s.User.UID {
			return fmt.Errorf("you are already requesting to join a duel. Leave the current request and try again")
		}
	}

	m.Lock()
	defer m.Unlock()

	r, ok := m.requests.Find(id)

	if !ok {
		return fmt.Errorf("the duel request you attempted to join does not exist")
	}

	if r.Guest.Some() {
		return fmt.Errorf("someone has already joined that duel")
	}

	r.Guest.Set(&MatchUser{
		ID:       s.User.UID,
		Username: s.User.Username,
		SocketID: s.UID,
	})

	return nil
}
