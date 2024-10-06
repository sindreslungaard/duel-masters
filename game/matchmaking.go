package game

import (
	"duel-masters/flags"
	"duel-masters/game/match"
	"duel-masters/internal"
	"duel-masters/server"
	"duel-masters/services"
	"fmt"
	"math/rand"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

var DefaultMatchNames = []string{
	"Kettou Da!",
	"I challenge you!",
	"Ikuzo!",
	"I'm ready!",
	"Koi!",
	"Bring it on!",
}

var Matchmaker = &Matchmaking{
	requests:    internal.NewConcurrentDictionary[MatchRequest](),
	broadcaster: func(msg interface{}) { logrus.Warn("Use of default matchmaking broadcaster") },
}

type MatchUser struct {
	ID       string
	Username string
	Color    string
	SocketID string
}

type MatchRequest struct {
	ID           string
	Host         MatchUser
	Guest        *internal.Option[MatchUser]
	Name         string
	Format       match.Format
	BlockedUsers internal.ConcurrentDictionary[bool]
	LinkCode     string
	EventUID     string
}

func (r *MatchRequest) Serialize() server.MatchRequestMessage {
	msg := server.MatchRequestMessage{
		ID:        r.ID,
		Name:      r.Name,
		HostID:    r.Host.ID,
		HostName:  r.Host.Username,
		HostColor: r.Host.Color,
		Format:    string(r.Format),
		LinkCode:  r.LinkCode,
		EventUID:  r.EventUID,
	}

	guest, ok := r.Guest.Unwrap()

	if !ok {
		return msg
	}

	msg.GuestID = guest.ID
	msg.GuestName = guest.Username
	msg.GuestColor = guest.Color

	return msg
}

type Matchmaking struct {
	sync.RWMutex
	requests    internal.ConcurrentDictionary[MatchRequest]
	broadcaster func(msg interface{})
	matchSystem *match.MatchSystem
}

func (m *Matchmaking) Initialize(f func(msg interface{}), sys *match.MatchSystem) {
	m.broadcaster = f
	m.matchSystem = sys
}

func (m *Matchmaking) NewRequest(s *server.Socket, name string, format match.Format, eventUID string) error {
	if !flags.NewMatchesEnabled {
		return fmt.Errorf("Match creation has been temporarily disabled")
	}

	m.Lock()
	defer m.Unlock()

	if len(name) > 50 {
		return fmt.Errorf("Please use a shorter name")
	}

	if name == "" {
		name = DefaultMatchNames[rand.Intn(len(DefaultMatchNames))]
	}

	if _, ok := m.requests.Find(s.User.UID); ok {
		return fmt.Errorf("You already have a duel request open")
	}

	for _, r := range m.requests.Iter() {
		if guest, ok := r.Guest.Unwrap(); ok && guest.ID == s.User.UID {
			return fmt.Errorf("Please leave the match request you are currently in before creating a new one")
		}
	}

	code, err := shortid.Generate()

	if err != nil {
		return fmt.Errorf("Failed to generate link code for match request")
	}

	if eventUID != "" {
		allowed, err := services.CanPlayerPlayEvent(s.User.UID, eventUID)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		if !allowed {
			return fmt.Errorf("Not allowed to play in this event")
		}
	}

	r := &MatchRequest{
		ID: s.User.UID,
		Host: MatchUser{
			ID:       s.User.UID,
			Username: s.User.Username,
			Color:    s.User.Color,
			SocketID: s.UID,
		},
		Guest:        internal.NewOption[MatchUser](),
		Name:         name,
		Format:       format,
		BlockedUsers: internal.NewConcurrentDictionary[bool](),
		LinkCode:     code,
		EventUID:     eventUID,
	}

	m.requests.Add(s.User.UID, r)
	m.BroadcastState()

	return nil
}

func (m *Matchmaking) Join(s *server.Socket, id string) error {
	m.Lock()
	defer m.Unlock()

	// check if guest is in another request
	for _, r := range m.requests.Iter() {
		if guest, ok := r.Guest.Unwrap(); ok && guest.ID == s.User.UID {
			return fmt.Errorf("You are already requesting to join a duel. Leave the current request and try again")
		}
	}

	_, ok := m.requests.Find(s.User.UID)
	if ok {
		return fmt.Errorf("Please close your open match request before joining someone elses")
	}

	r, ok := m.requests.Find(id)

	if !ok {
		return fmt.Errorf("The duel request you attempted to join does not exist")
	}

	if r.Host.ID == s.User.UID {
		return fmt.Errorf("You cannot join your own match")
	}

	if r.Guest.Some() {
		return fmt.Errorf("Someone has already joined that duel")
	}

	_, ok = r.BlockedUsers.Find(s.User.UID)
	if ok {
		return fmt.Errorf("You have been blocked from joining this match")
	}

	if r.EventUID != "" {
		allowed, err := services.CanPlayerPlayEvent(s.User.UID, r.EventUID)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		if !allowed {
			return fmt.Errorf("Not allowed to play in this event")
		}
	}

	r.Guest.Set(&MatchUser{
		ID:       s.User.UID,
		Username: s.User.Username,
		Color:    s.User.Color,
		SocketID: s.UID,
	})

	m.BroadcastState()

	sound := server.PlaySoundMessage{
		Header: "play_sound",
		Sound:  "request_accepted",
	}

	hostSocket, ok := server.Sockets.Find(r.Host.SocketID)

	if ok {
		hostSocket.Send(sound)
		chat(hostSocket, fmt.Sprintf("%s joined your duel request", s.User.Username))
	}

	s.Send(sound)

	return nil
}

func (m *Matchmaking) Leave(s *server.Socket) {
	m.Lock()
	defer m.Unlock()

	var broadcast bool

	for _, r := range m.requests.Iter() {
		guest, ok := r.Guest.Unwrap()

		if r.Host.ID == s.User.UID {
			m.requests.Remove(r.ID)

			if ok {
				guestSocket, _ := server.Sockets.Find(guest.SocketID)
				if guestSocket != nil {
					guestSocket.Warn("The match you requested to join has been closed by the host")
				}
			}

			broadcast = true
		} else if ok && guest.ID == s.User.UID {
			hostSocket, _ := server.Sockets.Find(r.Host.SocketID)
			if hostSocket != nil {
				chat(hostSocket, fmt.Sprintf("%s left your duel request", guest.Username))
			}
			r.Guest.Clear()
			broadcast = true
		}
	}

	if broadcast {
		m.BroadcastState()
	}
}

func (m *Matchmaking) Kick(s *server.Socket, requestId string, toKickId string) {
	m.Lock()
	defer m.Unlock()

	r, ok := m.requests.Find(requestId)

	if !ok {
		s.Warn("Failed to kick user from your duel request, it does not exist")
		return
	}

	if r.Host.ID != s.User.UID {
		s.Warn("Only the host can kick users")
		return
	}

	if len(r.BlockedUsers.Iter()) > 25 {
		s.Warn("You have reached the limit of how many users can be blocked from joining your match")
		return
	}

	r.BlockedUsers.Add(toKickId, nil)

	guest, ok := r.Guest.Unwrap()

	if !ok || guest.ID != toKickId {
		s.Warn("The user you tried to kick had already left, they are now blocked from joining again")
		return
	}

	r.Guest.Clear()
	guestSocket, ok := server.Sockets.Find(guest.SocketID)
	if ok {
		guestSocket.Warn("You were kicked from the match")
	}

	s.Warn(fmt.Sprintf("Successfully kicked %s. They are now blocked from joining again", guest.Username))
	m.BroadcastState()
}

func (m *Matchmaking) Start(s *server.Socket, requestId string) {
	m.Lock()
	defer m.Unlock()

	r, ok := m.requests.Find(requestId)

	if !ok {
		s.Warn("The match you attempted to start no longer exist")
		return
	}

	if r.Host.ID != s.User.UID {
		s.Warn("Only the host can start the match")
		return
	}

	guest, ok := r.Guest.Unwrap()

	if !ok {
		s.Warn("Can't start the match because there's no opponent")
		return
	}

	guestSocket, ok := server.Sockets.Find(guest.SocketID)

	if !ok {
		s.Warn("Failed to communicate with the opposing player, consider kicking them")
		return
	}

	match := m.matchSystem.NewMatch(r.Name, r.Host.ID, true, true, r.EventUID)

	msg := server.MatchForwardMessage{
		Header: "match_forward",
		ID:     match.ID,
	}

	s.Send(msg)
	guestSocket.Send(msg)

}

func (m *Matchmaking) Serialize() server.MatchReuestsListMessage {
	msg := server.MatchReuestsListMessage{
		Header:   "match_requests",
		Requests: []server.MatchRequestMessage{},
	}

	for _, r := range m.requests.Iter() {
		msg.Requests = append(msg.Requests, r.Serialize())
	}

	return msg
}

func (m *Matchmaking) BroadcastState() {
	m.broadcaster(m.Serialize())
}
