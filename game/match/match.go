package match

import (
	"context"
	"duel-masters/db"
	"duel-masters/server"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
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
		Turn:      1,
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

// CurrentPlayer returns either player1 or player2 based on who's turn it currently is
func (m *Match) CurrentPlayer() *PlayerReference {

	if m.Turn == 1 {
		return m.Player1
	}

	return m.Player2

}

// PlayerForSocket returns the player ref for a given socker or an error if the socket is not p1 or p2
func (m *Match) PlayerForSocket(s *server.Socket) (*PlayerReference, error) {

	if m.Player1.Socket == s {
		return m.Player1, nil
	}

	if m.Player2.Socket == s {
		return m.Player2, nil
	}

	return nil, errors.New("Socket is not player1 or player2")

}

// PlayerRef returns the player ref for a given player
func (m *Match) PlayerRef(p *Player) *PlayerReference {

	if m.Player1.Player == p {
		return m.Player1
	}

	return m.Player2

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

// BroadcastState sends the current game's state to both players, hiding the opponent's hand
func (m *Match) BroadcastState() {

	player1 := *m.Player1.Player.Denormalized()
	player2 := *m.Player2.Player.Denormalized()

	p1state := &server.MatchStateMessage{
		Header: "state_update",
		State: server.MatchState{
			MyTurn:       m.Turn == 1,
			HasAddedMana: m.Player1.Player.HasChargedMana,
			Me:           player1,
			Opponent:     player2,
		},
	}

	p2state := &server.MatchStateMessage{
		Header: "state_update",
		State: server.MatchState{
			MyTurn:       m.Turn == 2,
			HasAddedMana: m.Player2.Player.HasChargedMana,
			Me:           player2,
			Opponent:     player1,
		},
	}

	p1state.State.Opponent.Hand = make([]server.CardState, 0)
	p2state.State.Opponent.Hand = make([]server.CardState, 0)

	m.Player1.Socket.Send(p1state)
	m.Player2.Socket.Send(p2state)

}

// Warn sends a warning to the specified player ref
func Warn(p *PlayerReference, message string) {

	p.Socket.Send(server.WarningMessage{
		Header:  "warn",
		Message: message,
	})

}

// WarnPlayer sends a warning to the specified player
func (m *Match) WarnPlayer(p *Player, message string) {

	m.PlayerRef(p).Socket.Send(server.WarningMessage{
		Header:  "warn",
		Message: message,
	})

}

// HandleFx ...
func (m *Match) HandleFx(ctx *Context) {

	players := make([]*PlayerReference, 0)

	// The player in which turn it is is to be handled first
	if m.Turn == m.Player1.Player.Turn {
		players = append(players, m.Player1, m.Player2)
	} else {
		players = append(players, m.Player2, m.Player1)
	}

	cards := make([]*Card, 0)

	for _, p := range players {

		cards = append(cards, p.Player.Battlezone...)
		cards = append(cards, p.Player.Spellzone...)
		cards = append(cards, p.Player.Hand...)
		cards = append(cards, p.Player.Shieldzone...)

	}

	for _, card := range cards {

		for _, h := range card.handlers {

			if ctx.cancel {
				continue
			}

			h(card, ctx)

		}

	}

}

// NewAction prompts the user to make a selection of the specified []Cards
func (m *Match) NewAction(player *Player, cards []*Card, minSelections int, maxSelections int, text string, cancellable bool) {

	msg := &server.ActionMessage{
		Header:        "action",
		Cards:         denormalizeCards(cards),
		Text:          text,
		MinSelections: minSelections,
		MaxSelections: maxSelections,
		Cancellable:   cancellable,
	}

	m.PlayerRef(player).Socket.Send(msg)

}

// NewMultipartAction prompts the user to make a selection of the specified {string: []Cards}
func (m *Match) NewMultipartAction(player *Player, cards map[string][]*Card, minSelections int, maxSelections int, text string, cancellable bool) {

	cardMap := make(map[string][]server.CardState)

	for key, cards := range cards {
		cardMap[key] = denormalizeCards(cards)
	}

	msg := &server.MultipartActionMessage{
		Header:        "action",
		Cards:         cardMap,
		Text:          text,
		MinSelections: minSelections,
		MaxSelections: maxSelections,
		Cancellable:   cancellable,
	}

	m.PlayerRef(player).Socket.Send(msg)

}

// Start starts the match
func (m *Match) Start() {

	m.Player1.Player.ShuffleDeck()
	m.Player2.Player.ShuffleDeck()

	m.Player1.Player.InitShieldzone()
	m.Player2.Player.InitShieldzone()

	m.Player1.Player.DrawCards(5)
	m.Player2.Player.DrawCards(5)

	// match.turn is initialized as 1, so we only need to change it to 2
	// The opposite of what's defined here will start because BeginNewTurn() changes it
	if rand.Intn(100) >= 50 {
		m.Turn = 2
	}

	m.Chat("Server", "The duel has begun!")

	m.BeginNewTurn()

}

// BeginNewTurn starts a new turn
func (m *Match) BeginNewTurn() {

	if m.Turn == 1 {
		m.Turn = 2
	} else {
		m.Turn = 1
	}

	ctx := NewContext(m, &BeginTurnStep{})

	m.HandleFx(ctx)

	m.CurrentPlayer().Player.HasChargedMana = false

	m.BroadcastState()

	m.UntapStep()

}

// UntapStep ...
func (m *Match) UntapStep() {

	ctx := NewContext(m, &UntapStep{})

	m.HandleFx(ctx)

	m.StartOfTurnStep()

}

// StartOfTurnStep ...
func (m *Match) StartOfTurnStep() {

	ctx := NewContext(m, &StartOfTurnStep{})

	m.HandleFx(ctx)

	m.Chat("Server", fmt.Sprintf("Your turn, %s", m.CurrentPlayer().Socket.User.Username))

	m.DrawStep()

}

// DrawStep ...
func (m *Match) DrawStep() {

	ctx := NewContext(m, &DrawStep{})

	m.HandleFx(ctx)

	m.CurrentPlayer().Player.DrawCards(1)

	m.BroadcastState()

	m.Chat("Server", fmt.Sprintf("%s drew 1 card", m.CurrentPlayer().Socket.User.Username))

	m.ChargeStep()

}

// ChargeStep ...
func (m *Match) ChargeStep() {

	ctx := NewContext(m, &ChargeStep{})

	m.HandleFx(ctx)

}

// EndStep ...
func (m *Match) EndStep() {

	ctx := NewContext(m, &EndStep{})

	m.HandleFx(ctx)

	m.Chat("Server", fmt.Sprintf("%s ended their turn", m.CurrentPlayer().Socket.User.Username))

	m.EndOfTurnTriggers()

}

// EndOfTurnTriggers ...
func (m *Match) EndOfTurnTriggers() {

	ctx := NewContext(m, &EndOfTurnStep{})

	m.HandleFx(ctx)

	m.BeginNewTurn()

}

// EndTurn is called when the player attempts to end their turn
// If the context is not cancelled by a card, the EndStep is called
func (m *Match) EndTurn() {

	ctx := NewContext(m, &EndTurnEvent{})

	m.HandleFx(ctx)

	if !ctx.cancel {
		m.EndStep()
	}

}

// ChargeMana is called when the player attempts to charge mana
func (m *Match) ChargeMana(p *PlayerReference, cardID string) {

	if p.Player.HasChargedMana {
		Warn(p, "You have already charged mana this round")
		return
	}

	if card, err := p.Player.MoveCard(cardID, HAND, MANAZONE); err == nil {
		p.Player.HasChargedMana = true
		m.BroadcastState()
		m.Chat("Server", fmt.Sprintf("%s was added to %s's manazone", card.Name, p.Socket.User.Username))
	}

}

// PlayCard is called when the player attempts to play a card
func (m *Match) PlayCard(p *PlayerReference, cardID string) {

	m.HandleFx(NewContext(m, &PlayCardEvent{
		CardID: cardID,
	}))

	m.BroadcastState()

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

	case "choose_deck":
		{

			p, err := m.PlayerForSocket(s)

			if err != nil {
				return
			}

			var msg struct {
				UID string `json:"uid"`
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			var deck db.Deck

			if err := db.Collection("decks").FindOne(context.TODO(), bson.M{"uid": msg.UID}).Decode(&deck); err != nil {
				return
			}

			p.Player.CreateDeck(deck.Cards)

			m.Chat("Server", fmt.Sprintf("%s has chosen their deck", s.User.Username))

			p.Player.Ready = true

			if m.Player1.Player.Ready && m.Player2.Player.Ready {
				m.Start()
			}

		}

	case "add_to_manazone":
		{

			p, err := m.PlayerForSocket(s)

			if err != nil {
				return
			}

			if m.Turn != p.Player.Turn {
				return
			}

			var msg struct {
				ID string `json:"virtualId"`
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			m.ChargeMana(p, msg.ID)

		}

	case "end_turn":
		{

			p, err := m.PlayerForSocket(s)

			if err != nil {
				return
			}

			if m.Turn != p.Player.Turn {
				return
			}

			m.EndTurn()

		}

	case "add_to_playzone":
		{

			p, err := m.PlayerForSocket(s)

			if err != nil {
				return
			}

			if m.Turn != p.Player.Turn {
				return
			}

			var msg struct {
				ID string `json:"virtualId"`
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			m.PlayCard(p, msg.ID)

		}

	case "action":
		{

			p, err := m.PlayerForSocket(s)

			if err != nil {
				return
			}

			var msg PlayerAction

			if err := json.Unmarshal(data, &msg); err != nil {
				Warn(p, "Invalid selection")
				return
			}

			p.Player.Action <- msg

		}

	default:
		{
			logrus.Debugf("Received message in incorrect format: %v", string(data))
		}

	}

}
