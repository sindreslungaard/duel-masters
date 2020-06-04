package game

import (
	"duel-masters/game/match"
	"duel-masters/server"
	"encoding/json"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	messageBufferSize int = 100
)

var messages = make([]server.LobbyChatMessage, 0)
var messagesMutex = &sync.Mutex{}

var subscribers = make([]*server.Socket, 0)
var subscribersMutex = &sync.Mutex{}

var userCache server.UserListMessage = server.GetUserList()
var matchCache server.MatchesListMessage = server.MatchesListMessage{
	Header:  "matches",
	Matches: make([]server.MatchMessage, 0),
}

var lobby = &Lobby{}

// Lobby struct is used to create a Hub that can parse messages from the websocket server
type Lobby struct{}

// Name just returns "lobby", obligatory for a hub
func (l *Lobby) Name() string {
	return "lobby"
}

// GetLobby returns a reference to the lobby
func GetLobby() *Lobby {
	return lobby
}

// Broadcast sends a message to all subscribed sockets
func Broadcast(msg interface{}) {
	subscribersMutex.Lock()
	defer subscribersMutex.Unlock()

	for _, subscriber := range subscribers {
		go subscriber.Send(msg)
	}
}

// StartTicker starts the lobby ticker
func (l *Lobby) StartTicker() {

	ticker := time.NewTicker(30 * time.Second) // tick every 30 seconds

	defer ticker.Stop()

	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered from lobby ticker. %v", r)
		}
	}()

	go ListenForMatchListUpdates()

	for {

		select {
		case <-ticker.C:
			{
				UpdateUserCache()
				Broadcast(userCache)
			}
		}

	}
}

func UpdateUserCache() {
	userCache = server.GetUserList()
}

// ListenForMatchListUpdates broadcasts changes to the open matches to all lobby subscribers
func ListenForMatchListUpdates() {

	for {

		update := <-match.LobbyMatchList()

		matchCache = update

		Broadcast(update)

	}

}

// Parse websocket messages
func (l *Lobby) Parse(s *server.Socket, data []byte) {

	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from parsing a message in lobby. %v", r)
		}
	}()

	var message server.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return
	}

	switch message.Header {

	case "subscribe":
		{
			subscribersMutex.Lock()
			defer subscribersMutex.Unlock()

			for _, subscriber := range subscribers {
				if subscriber == s {
					return
				}
			}

			subscribers = append(subscribers, s)

			// Send chat messages
			s.Send(server.LobbyChatMessages{
				Header:   "chat",
				Messages: messages,
			})

			// Update and send user list
			UpdateUserCache()
			s.Send(userCache)

			// Send match list

			s.Send(matchCache)

		}

	case "chat":
		{

			var msg struct {
				Message string `json:"message"`
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			messagesMutex.Lock()
			defer messagesMutex.Unlock()

			if len(messages) >= messageBufferSize {
				_, messages = messages[0], messages[1:]
			}

			chatMsg := server.LobbyChatMessage{
				Username:  s.User.Username,
				Color:     s.User.Color,
				Message:   msg.Message,
				Timestamp: int(time.Now().Unix()),
			}

			toBroadcast := server.LobbyChatMessages{
				Header:   "chat",
				Messages: []server.LobbyChatMessage{chatMsg},
			}

			messages = append(messages, chatMsg)

			Broadcast(toBroadcast)

		}

	}

}

// OnSocketClose is called when a socket disconnects
func (l *Lobby) OnSocketClose(s *server.Socket) {

	subscribersMutex.Lock()
	defer subscribersMutex.Unlock()

	subscribersUpdate := make([]*server.Socket, 0)

	for _, subscriber := range subscribers {

		if subscriber != s {
			subscribersUpdate = append(subscribersUpdate, subscriber)
		}

	}

	subscribers = subscribersUpdate

}
