package match

import (
	"context"
	"duel-masters/db"
	"duel-masters/server"
	"encoding/json"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
	"go.mongodb.org/mongo-driver/bson"
)

var matches = make(map[string]*Match)
var matchesMutex = sync.Mutex{}

// Match struct
type Match struct {
	ID        string           `json:"id"`
	MatchName string           `json:"name"`
	HostID    string           `json:"-"`
	Player1   *PlayerReference `json:"-"`
	Player2   *PlayerReference `json:"-"`
	Turn      byte             `json:"-"`
	Started   bool             `json:"started"`
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
		Started:   false,
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
	return m.Turn == p.Turn
}

// ColorChat sends a chat message with color
func (m *Match) ColorChat(sender string, message string, color string) {
	msg := &server.ChatMessage{
		Header:  "chat",
		Message: message,
		Sender:  sender,
		Color:   color,
	}

	m.Player1.Socket.Send(msg)
	m.Player2.Socket.Send(msg)
}

// Chat sends a chat message with the default color
func (m *Match) Chat(sender string, message string) {
	m.ColorChat(sender, message, "#ccc")
}

// Parse handles websocket messages in this Hub
func (m *Match) Parse(s *server.Socket, data []byte) {

	var message server.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return
	}

	switch message.Header {

	case "join_match":
		{

			// TODO: spectators?
			if m.Started {
				return
			}

			// This is player1
			if s.User.UID == m.HostID {

				if m.Player1 != nil {
					// TODO: Allow reconnect?
					logrus.Debug("Attempt to join as Player1 multiple times")
					return
				}

				p := NewPlayer(1)

				m.Player1 = NewPlayerReference(p, s)

			}

			// This is player2
			if s.User.UID != m.HostID {

				if m.Player2 != nil {
					// TODO: Allow reconnect?
					logrus.Debug("Attempt to join as Player2 multiple times")
					return
				}

				p := NewPlayer(2)

				m.Player2 = NewPlayerReference(p, s)

			}

			// If both players have joined, prompt them to choose their decks
			if m.Player1 != nil && m.Player2 != nil {

				collection := db.Collection("decks")

				cur, err := collection.Find(context.TODO(), bson.M{
					"$or": []bson.M{
						{"owner": m.Player1.Socket.User.UID},
						{"owner": m.Player2.Socket.User.UID},
						{"standard": true},
					},
				})

				if err != nil {
					logrus.Error(err)
					return
				}

				defer cur.Close(context.TODO())

				player1decks := make([]db.Deck, 0)
				player2decks := make([]db.Deck, 0)

				for cur.Next(context.TODO()) {

					var deck db.Deck

					if err := cur.Decode(&deck); err != nil {
						continue
					}

					if deck.Owner == m.Player1.Socket.User.UID || deck.Standard {
						player1decks = append(player1decks, deck)
					}

					if deck.Owner == m.Player2.Socket.User.UID || deck.Standard {
						player2decks = append(player2decks, deck)
					}

				}

				m.Player1.Socket.Send(server.DecksMessage{
					Header: "choose_deck",
					Decks:  player1decks,
				})

				m.Player2.Socket.Send(server.DecksMessage{
					Header: "choose_deck",
					Decks:  player2decks,
				})

				m.Chat("Server", "Waiting for both players to choose a deck")

			}

		}

	case "chat":
		{

			// Allow other sockets than player1 and player2 to chat?

			var msg struct {
				Message string `json:"message"`
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			m.ColorChat(s.User.Username, msg.Message, "#79dced")
		}

	}

}
