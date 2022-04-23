package match

import (
	"context"
	"duel-masters/db"
	"duel-masters/game/cnd"
	"duel-masters/server"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"runtime/debug"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
	"go.mongodb.org/mongo-driver/bson"
)

var matches = make(map[string]*Match)
var matchesMutex = sync.Mutex{}

// Get returns a *Match from the specified id
func Get(id string) (*Match, error) {
	matchesMutex.Lock()
	defer matchesMutex.Unlock()

	if m, ok := matches[id]; ok {
		return m, nil
	}

	return nil, errors.New("Not found")
}

var lobbyMatches = make(chan server.MatchesListMessage)

// Match struct
type Match struct {
	ID                string           `json:"id"`
	MatchName         string           `json:"name"`
	HostID            string           `json:"-"`
	Player1           *PlayerReference `json:"-"`
	Player2           *PlayerReference `json:"-"`
	spectators        Spectators       `json:"-"`
	persistentEffects map[int]PersistentEffect
	Turn              byte `json:"-"`
	Started           bool `json:"started"`
	Visible           bool `json:"visible"`
	Step              interface{}

	created     int64
	ending      bool
	closed      bool
	isFirstTurn bool

	quit chan bool
}

// Matches returns a list of the current matches
func Matches() []string {
	result := make([]string, 0)
	matchesMutex.Lock()
	defer matchesMutex.Unlock()
	for id := range matches {
		result = append(result, id)
	}
	return result
}

// New returns a new match object
func New(matchName string, hostID string, visible bool) *Match {

	id, err := shortid.Generate()

	if err != nil {
		id = uuid.New().String()
	}

	m := &Match{
		ID:                id,
		MatchName:         matchName,
		HostID:            hostID,
		spectators:        Spectators{users: map[string]Spectator{}},
		persistentEffects: make(map[int]PersistentEffect),
		Turn:              1,
		Started:           false,
		Visible:           visible,

		created:     time.Now().Unix(),
		ending:      false,
		isFirstTurn: true,

		quit: make(chan bool),
	}

	matchesMutex.Lock()

	matches[id] = m

	matchesMutex.Unlock()

	go m.startTicker()

	logrus.Debugf("Created match %s", id)

	return m

}

// Name just returns "match", obligatory for a hub
func (m *Match) Name() string {
	return "match"
}

// LobbyMatchList returns the channel to receive match list updates
func LobbyMatchList() chan server.MatchesListMessage {
	return lobbyMatches
}

// UpdateMatchList sends a server.MatchesListMessage through the lobby channel
func UpdateMatchList() {

	matchesMutex.Lock()
	defer matchesMutex.Unlock()

	matchesMessage := make([]server.MatchMessage, 0)

	for _, match := range matches {

		if !match.Visible {
			continue
		}

		if match.ending {
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

	update := server.MatchesListMessage{
		Header:  "matches",
		Matches: matchesMessage,
	}

	lobbyMatches <- update

}

func (m *Match) startTicker() {

	ticker := time.NewTicker(10 * time.Second) // tick every 10 seconds

	defer ticker.Stop()
	defer m.Dispose()
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from match ticker. %v", r)
			debug.PrintStack()
		}
	}()

	for {

		select {
		case <-m.quit:
			{
				logrus.Debugf("Closing match %s", m.ID)
				m.ending = true
				return
			}
		case <-ticker.C:
			{

				// Close the match if it was not started within 10 minutes of creation
				if !m.Started && m.created < time.Now().Unix()-60*10 {
					logrus.Debugf("Closing match %s", m.ID)
					return
				}

			}
		}

	}
}

// Dispose closes the match, disconnects the clients and removes all references to it
func (m *Match) Dispose() {

	if m.closed {
		return
	}

	m.closed = true

	logrus.Debugf("Disposing match %s", m.ID)

	defer func() {
		if r := recover(); r != nil {
			logrus.Warningf("Recovered from disposing a match. %v", r)
			debug.PrintStack()
		}
	}()

	m.spectators.Lock()
	defer m.spectators.Unlock()
	for _, spectator := range m.spectators.users {
		if spectator.Socket == nil {
			continue
		}
		spectator.Socket.Close()
		spectator.Socket = nil
	}

	if m.Player1 != nil {
		m.Player1.Socket.Close()
		m.Player1.Player.Dispose()
	}

	if m.Player2 != nil {
		m.Player2.Socket.Close()
		m.Player2.Player.Dispose()
	}

	matchesMutex.Lock()

	close(m.quit)

	delete(matches, m.ID)

	matchesMutex.Unlock()

	logrus.Debugf("Closed match with id %s", m.ID)

	UpdateMatchList()

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

// Opponent returns the opponent of the given player
func (m *Match) Opponent(p *Player) *Player {
	if m.Player1.Player == p {
		return m.Player2.Player
	}

	return m.Player1.Player
}

// GetPower returns the power of a given card after applying conditions
func (m *Match) GetPower(card *Card, isAttacking bool) int {

	power := card.Power

	for _, condition := range card.Conditions() {

		switch condition.ID {

		case cnd.PowerAmplifier:
			{
				if val, ok := condition.Val.(int); ok {
					power += val
				}
			}

		case cnd.PowerAttacker:
			{
				if !isAttacking {
					continue
				}

				if val, ok := condition.Val.(int); ok {
					power += val
				}
			}

		}

	}

	power += card.PowerModifier(m, isAttacking)

	e := &GetPowerEvent{
		Card:      card,
		Attacking: isAttacking,
		Power:     power,
	}
	m.HandleFx(NewContext(m, e))

	return e.Power

}

// CastSpell Fires a SpellCast event
func (m *Match) CastSpell(card *Card, fromShield bool) {

	m.HandleFx(NewContext(m, &SpellCast{
		CardID:     card.ID,
		FromShield: fromShield,
	}))

	m.BroadcastState()

}

// Battle handles a battle between two creatures
func (m *Match) Battle(attacker *Card, defender *Card, blocked bool) {

	attackerPower := m.GetPower(attacker, true)
	defenderPower := m.GetPower(defender, false)

	m.HandleFx(NewContext(m, &AttackConfirmed{CardID: attacker.ID, Player: false, Creature: true}))
	m.HandleFx(NewContext(m, &Battle{Attacker: attacker, Defender: defender, Blocked: blocked}))

	if attackerPower > defenderPower {
		m.HandleFx(NewContext(m, &CreatureDestroyed{Card: defender, Source: attacker, Blocked: blocked}))
		m.Chat("Server", fmt.Sprintf("%s (%v) was destroyed by %s (%v)", defender.Name, defenderPower, attacker.Name, attackerPower))
	} else if attackerPower == defenderPower {
		m.HandleFx(NewContext(m, &CreatureDestroyed{Card: attacker, Source: defender, Blocked: blocked}))
		m.Chat("Server", fmt.Sprintf("%s (%v) was destroyed by %s (%v)", attacker.Name, attackerPower, defender.Name, defenderPower))
		m.HandleFx(NewContext(m, &CreatureDestroyed{Card: defender, Source: attacker, Blocked: blocked}))
		m.Chat("Server", fmt.Sprintf("%s (%v) was destroyed by %s (%v)", defender.Name, defenderPower, attacker.Name, attackerPower))
	} else if attackerPower < defenderPower {
		m.HandleFx(NewContext(m, &CreatureDestroyed{Card: attacker, Source: defender, Blocked: blocked}))
		m.Chat("Server", fmt.Sprintf("%s (%v) was destroyed by %s (%v)", attacker.Name, attackerPower, defender.Name, defenderPower))
	}

	m.BroadcastState()

}

// Destroy sends the given card to its players graveyard
func (m *Match) Destroy(card *Card, source *Card, context CreatureDestroyedContext) {

	m.HandleFx(NewContext(m, &CreatureDestroyed{Card: card, Source: source, Context: context}))
	m.Chat("Server", fmt.Sprintf("%s (%v) was destroyed by %s", card.Name, m.GetPower(card, false), source.Name))

}

// MoveCard moves a card and sends a chat message about what source moved it
func (m *Match) MoveCard(card *Card, destination string, source *Card) {

	_, err := card.Player.MoveCard(card.ID, card.Zone, destination)

	if err != nil {
		return
	}

	m.Chat("Server", fmt.Sprintf("%s was moved to %s %s by %s", card.Name, card.Player.Username(), destination, source.Name))

}

// MoveCardToFront moves a card and sends a chat message about what source moved it
func (m *Match) MoveCardToFront(card *Card, destination string, source *Card) {

	_, err := card.Player.MoveCardToFront(card.ID, card.Zone, destination)

	if err != nil {
		return
	}

	m.Chat("Server", fmt.Sprintf("%s was moved to %s's %s by %s", card.Name, card.Player.Username(), destination, source.Name))

}

// BreakShields breaks the given shields and handles shieldtriggers
func (m *Match) BreakShields(shields []*Card) {

	if len(shields) < 1 {
		return
	}

	m.Chat("Server", fmt.Sprintf("%v of %v's shields were broken", len(shields), m.PlayerRef(shields[0].Player).Socket.User.Username))

	for _, shield := range shields {

		card, err := shield.Player.MoveCard(shield.ID, SHIELDZONE, HAND)

		if err != nil {
			continue
		}

		// Handle shield triggers
		if card.HasCondition(cnd.ShieldTrigger) {

			ctx := NewContext(m, &ShieldTriggerEvent{
				Card: card,
			})

			m.HandleFx(ctx)

			if ctx.Cancelled() {
				continue
			}

			m.Wait(m.Opponent(card.Player), "Waiting for your opponent to make an action")

			m.NewAction(card.Player, []*Card{card}, 1, 1, "Shield trigger! Choose the card to use for free or close to keep it in your hand", true)

			for {

				action := <-card.Player.Action

				if action.Cancel {
					m.CloseAction(card.Player)
					break
				}

				if len(action.Cards) < 1 {
					m.DefaultActionWarning(card.Player)
					continue
				}

				if card.HasCondition(cnd.Spell) {
					m.CastSpell(card, true)
				} else {
					m.MoveCard(card, BATTLEZONE, card)
				}

				m.CloseAction(card.Player)

				break

			}

			m.EndWait(m.Opponent(card.Player))

		}

	}

}

// End ends the match
func (m *Match) End(winner *Player, winnerStr string) {

	logrus.Debugf("Attempting to end match %s", m.ID)

	if m.ending {
		logrus.Debugf("Cannot end match, %s is already ending", m.ID)
		return
	}

	if m.Started {

		m.Broadcast(server.WarningMessage{
			Header:  "error",
			Message: winnerStr,
		})

	}

	m.quit <- true

}

// ColorChat sends a chat message with color
func (m *Match) ColorChat(sender string, message string, color string) {
	msg := &server.ChatMessage{
		Header:  "chat",
		Message: message,
		Sender:  sender,
		Color:   color,
	}

	m.Broadcast(msg)
}

// Chat sends a chat message with the default color
func (m *Match) Chat(sender string, message string) {
	m.ColorChat(sender, message, "#ccc")
}

func (m *Match) Broadcast(msg interface{}) {

	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered during Broadcast(). %v", r)
		}
	}()

	m.Player1.Socket.Send(msg)
	m.Player2.Socket.Send(msg)

	// could fail due to concurrent updates to spectators map
	// but don't want to lock it here as broadcast will be used everywhere
	// so recover will deal with those special occasians where it fails
	for _, spectator := range m.spectators.users {
		if spectator.Socket == nil {
			continue
		}
		spectator.Socket.Send(msg)
	}

}

// BroadcastState sends the current game's state to both players, hiding the opponent's hand
func (m *Match) BroadcastState() {

	player1 := *m.Player1.Player.Denormalized()
	player1.Username = m.Player1.Username
	player1.Color = m.Player1.Color

	player2 := *m.Player2.Player.Denormalized()
	player2.Username = m.Player2.Username
	player2.Color = m.Player2.Color

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

	spectatorState := &server.MatchStateMessage{
		Header: "state_update",
		State: server.MatchState{
			MyTurn:       false,
			HasAddedMana: false,
			Me:           player1,
			Opponent:     player2,
			Spectator:    true,
		},
	}

	p1state.State.Opponent.Hand = make([]server.CardState, 0)
	p2state.State.Opponent.Hand = make([]server.CardState, 0)
	spectatorState.State.Me.Hand = make([]server.CardState, 0)
	spectatorState.State.Opponent.Hand = make([]server.CardState, 0)

	m.Player1.Socket.Send(p1state)
	m.Player2.Socket.Send(p2state)

	m.spectators.RLock()
	defer m.spectators.RUnlock()

	for _, spectator := range m.spectators.users {
		if spectator.Socket == nil {
			continue
		}
		spectator.Socket.Send(spectatorState)
	}

}

// Warn sends a warning to the specified player ref
func Warn(p *PlayerReference, message string) {

	p.Socket.Send(server.WarningMessage{
		Header:  "warn",
		Message: message,
	})

}

// WarnError sends an error message to the specified player ref
func WarnError(p *PlayerReference, message string) {

	p.Socket.Send(server.WarningMessage{
		Header:  "error",
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

// ActionWarning adds an error message to the players current action popup
func (m *Match) ActionWarning(p *Player, message string) {
	m.PlayerRef(p).Socket.Send(server.ActionWarningMessage{
		Header:  "action_error",
		Message: message,
	})
}

// DefaultActionWarning sends an actionw arning with a predefined message
func (m *Match) DefaultActionWarning(p *Player) {
	m.PlayerRef(p).Socket.Send(server.ActionWarningMessage{
		Header:  "action_error",
		Message: "Your selection of cards does not fulfill the requirements",
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

		cards = append(cards, p.Player.battlezone...)
		cards = append(cards, p.Player.spellzone...)
		cards = append(cards, p.Player.hand...)
		cards = append(cards, p.Player.shieldzone...)
		cards = append(cards, p.Player.hiddenzone...)
		cards = append(cards, p.Player.manazone...)
		cards = append(cards, p.Player.graveyard...)
		cards = append(cards, p.Player.deck...)

	}

	// Handle persistent effects
	for _, fx := range m.persistentEffects {
		fx.effect(ctx, fx.exit)
	}

	// Handle regular card effects c.Use(...
	for _, card := range cards {

		for _, h := range card.handlers {

			if ctx.cancel {
				return
			}

			h(card, ctx)

		}

	}

	// Handle ctx.ScheduleAfter effects
	for _, h := range ctx.postFxs {

		if ctx.cancel {
			return
		}

		h()

	}

}

// NewAction prompts the user to make a selection of the specified []Cards
func (m *Match) NewAction(player *Player, cards []*Card, minSelections int, maxSelections int, text string, cancellable bool) {

	msg := &server.ActionMessage{
		Header:        "action",
		Cards:         denormalizeCards(cards, false),
		Text:          text,
		MinSelections: minSelections,
		MaxSelections: maxSelections,
		Cancellable:   cancellable,
	}

	m.PlayerRef(player).Socket.Send(msg)

}

// NewBacksideAction prompts the user to make a selection of the specified cards without their names or images
func (m *Match) NewBacksideAction(player *Player, cards []*Card, minSelections int, maxSelections int, text string, cancellable bool) {

	msg := &server.ActionMessage{
		Header:        "action",
		Cards:         denormalizeCards(cards, true),
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
		cardMap[key] = denormalizeCards(cards, false)
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

// CloseAction closes the card selection popup for the given player
func (m *Match) CloseAction(p *Player) {
	m.PlayerRef(p).Socket.Send(server.Message{
		Header: "close_action",
	})
}

// Wait sends a waiting popup with a message to the specified player
func (m *Match) Wait(p *Player, message string) {
	m.PlayerRef(p).Socket.Send(server.WaitMessage{
		Header:  "wait",
		Message: message,
	})
}

// EndWait closes the waiting popup for the specified player
func (m *Match) EndWait(p *Player) {
	m.PlayerRef(p).Socket.Send(server.Message{
		Header: "end_wait",
	})
}

// ShowCards shows the specified cards to the player with a message of why it is being shown
func (m *Match) ShowCards(p *Player, message string, cards []string) {
	m.PlayerRef(p).Socket.Send(server.ShowCardsMessage{
		Header:  "show_cards",
		Message: message,
		Cards:   cards,
	})
}

// Start starts the match
func (m *Match) Start() {

	m.Started = true

	UpdateMatchList()

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

	m.Step = &BeginTurnStep{}

	if m.Turn == 1 {
		m.Turn = 2
	} else {
		m.Turn = 1
	}

	ctx := NewContext(m, m.Step)

	m.HandleFx(ctx)

	m.CurrentPlayer().Player.HasChargedMana = false
	m.CurrentPlayer().Player.CanChargeMana = true

	m.BroadcastState()

	m.UntapStep()

}

// UntapStep ...
func (m *Match) UntapStep() {

	m.Step = &UntapStep{}

	if mana, err := m.CurrentPlayer().Player.Container(MANAZONE); err == nil {
		for _, c := range mana {
			c.Tapped = false
		}
	}

	ctx := NewContext(m, m.Step)

	m.HandleFx(ctx)

	m.StartOfTurnStep()

}

// StartOfTurnStep ...
func (m *Match) StartOfTurnStep() {

	m.Step = &StartOfTurnStep{}

	ctx := NewContext(m, m.Step)

	m.HandleFx(ctx)

	m.Chat("Server", fmt.Sprintf("Your turn, %s", m.CurrentPlayer().Socket.User.Username))

	m.DrawStep()

}

// DrawStep ...
func (m *Match) DrawStep() {

	m.Step = &DrawStep{}

	ctx := NewContext(m, m.Step)

	m.HandleFx(ctx)

	if m.isFirstTurn == false {
		m.CurrentPlayer().Player.DrawCards(1)
	}

	m.BroadcastState()

	m.ChargeStep()

}

// ChargeStep ...
func (m *Match) ChargeStep() {

	m.Step = &ChargeStep{}

	ctx := NewContext(m, m.Step)

	m.HandleFx(ctx)

}

// EndStep ...
func (m *Match) EndStep() {

	m.Step = &EndStep{}

	ctx := NewContext(m, m.Step)

	m.HandleFx(ctx)

	m.Chat("Server", fmt.Sprintf("%s ended their turn", m.CurrentPlayer().Socket.User.Username))

	m.EndOfTurnTriggers()

}

// EndOfTurnTriggers ...
func (m *Match) EndOfTurnTriggers() {

	m.Step = &EndOfTurnStep{}

	if cards, err := m.CurrentPlayer().Player.Container(BATTLEZONE); err == nil {
		for _, c := range cards {
			c.ClearConditions()
		}
	}

	m.isFirstTurn = false

	ctx := NewContext(m, m.Step)

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

	if !p.Player.CanChargeMana {
		Warn(p, "You can't charge mana after playing or attacking with creatures/spells")
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

	ctx := NewContext(m, &PlayCardEvent{
		CardID: cardID,
	})

	m.HandleFx(ctx)

	if !ctx.Cancelled() {
		if _, ok := m.Step.(*MainStep); !ok {
			m.Step = &MainStep{}
		}

		p.Player.CanChargeMana = false
	}

	m.BroadcastState()

}

// AttackPlayer is called when the player attempts to attack the opposing player
func (m *Match) AttackPlayer(p *PlayerReference, cardID string) {

	_, err := p.Player.GetCard(cardID, BATTLEZONE)

	if err != nil {
		Warn(p, "The creature you tried to attack with is not in the battlezone")
		return
	}

	ctx := NewContext(m, &AttackPlayer{
		CardID:   cardID,
		Blockers: make([]*Card, 0),
	})

	m.HandleFx(ctx)

	if !ctx.Cancelled() {
		if _, ok := m.Step.(*AttackStep); !ok {
			m.Step = &AttackStep{}
		}

		p.Player.CanChargeMana = false
	}

	m.BroadcastState()

}

// AttackCreature is called when the player attempts to attack the opposing player
func (m *Match) AttackCreature(p *PlayerReference, cardID string) {

	_, err := p.Player.GetCard(cardID, BATTLEZONE)

	if err != nil {
		Warn(p, "The creature you tried to attack with is not in the battlezone")
		return
	}

	ctx := NewContext(m, &AttackCreature{
		CardID:              cardID,
		Blockers:            make([]*Card, 0),
		AttackableCreatures: make([]*Card, 0),
	})

	m.HandleFx(ctx)

	if !ctx.Cancelled() {
		if _, ok := m.Step.(*AttackStep); !ok {
			m.Step = &AttackStep{}
		}

		p.Player.CanChargeMana = false
	}

	m.BroadcastState()

}

// Parse handles websocket messages in this Hub
func (m *Match) Parse(s *server.Socket, data []byte) {

	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered after parsing message in match. %v", r)
			debug.PrintStack()
		}
	}()

	var message server.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return
	}

	switch message.Header {

	case "mpong":
		{

			p, err := m.PlayerForSocket(s)

			if err != nil {
				return
			}

			p.LastPong = time.Now().Unix()

		}

	case "join_match":
		{

			// player reconnect
			if m.Started {
				// p1 attempting to reconnect
				if m.Player1 != nil && m.Player1.UID == s.User.UID {

					if m.Player1.Socket != nil {
						m.Player1.Socket.Close()
					}

					m.Player1.Socket = s

					if m.Player2 != nil && m.Player2.Socket != nil {
						m.Player2.Socket.Send(server.Message{
							Header: "opponent_reconnected",
						})
					}

					m.BroadcastState()
					m.Chat("Server", s.User.Username+" reconnected")

					return

					// p2 attempting to reconnect
				} else if m.Player2 != nil && m.Player2.UID == s.User.UID {

					if m.Player2.Socket != nil {
						m.Player2.Socket.Close()
					}

					m.Player2.Socket = s

					if m.Player1 != nil && m.Player1.Socket != nil {
						m.Player1.Socket.Send(server.Message{
							Header: "opponent_reconnected",
						})
					}

					m.BroadcastState()
					m.Chat("Server", s.User.Username+" reconnected")

					return

				} else {
					// Spectators

					m.spectators.RLock()
					spectator, ok := m.spectators.users[s.User.UID]
					m.spectators.RUnlock()

					// this user is already spectating, swap connection to new one
					if ok {
						spectator.Socket.Send(server.WarningMessage{
							Header:  "error",
							Message: "You started spectating from a new connection, closing this one...",
						})
						// this removes the existing spectator from the match
						spectator.Socket.Close()
					}

					m.spectators.Lock()
					m.spectators.users[s.User.UID] = Spectator{
						UID:      s.User.UID,
						Username: s.User.Username,
						Color:    s.User.Color,
						Socket:   s,
						LastPong: time.Now().Unix(),
					}
					m.spectators.Unlock()

					m.Chat("Server", fmt.Sprintf("%s started spectating", s.User.Username))
					m.BroadcastState()
					return
				}
			}

			// This is player1
			if s.User.UID == m.HostID {

				if m.Player1 != nil {
					logrus.Debug("Attempt to join as Player1 multiple times")
					s.Send(server.WarningMessage{
						Header:  "error",
						Message: "You have already joined this match",
					})
					s.Close()
					return
				}

				p := NewPlayer(m, 1)

				m.Player1 = NewPlayerReference(p, s)

			}

			// This is player2
			if s.User.UID != m.HostID {

				if m.Player2 != nil {
					logrus.Debug("Attempt to join as Player2 multiple times")
					s.Send(server.WarningMessage{
						Header:  "error",
						Message: "This match has already started, you cannot join it",
					})
					s.Close()
					return
				}

				p := NewPlayer(m, 2)

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

			UpdateMatchList()

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

			runes := []rune(msg.Message)
			if string(runes[0:4]) == "/add" {

				hasRights := false

				for _, permission := range s.User.Permissions {
					if permission == "admin" {
						hasRights = true
					}
				}

				if !hasRights {
					return
				}

				// Spawn card
				m.CurrentPlayer().Player.SpawnCard(string(runes[5:]))
				m.BroadcastState()

				return
			}

			m.ColorChat(s.User.Username, msg.Message, s.User.Color)
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

			msg := PlayerAction{}

			if err := json.Unmarshal(data, &msg); err != nil {
				m.ActionWarning(p.Player, "Invalid selection")
				return
			}

			// Check to see if the client is trying something fishy with selecting the same card multiple times
			for _, c := range msg.Cards {
				count := 0
				for _, c2 := range msg.Cards {
					if c == c2 {
						count++
					}
				}
				if count >= 2 {
					m.ActionWarning(p.Player, "You cannot select the same card multiple times")
					return
				}
			}

			p.Player.Action <- msg

		}

	case "attack_player":
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

			m.AttackPlayer(p, msg.ID)

		}

	case "attack_creature":
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

			m.AttackCreature(p, msg.ID)

		}

	default:
		{
			logrus.Debugf("Received message in incorrect format: %v", string(data))
		}

	}

}

// OnSocketClose is called when a socket disconnects
func (m *Match) OnSocketClose(s *server.Socket) {

	if m.closed {
		return
	}

	if !m.Started {
		m.quit <- true
		return
	}

	// is this a spectator leaving?
	m.spectators.Lock()
	defer m.spectators.Unlock()
	spectator, ok := m.spectators.users[s.User.UID]

	if ok {
		m.Chat("Server", fmt.Sprintf("%s stopped spectating", spectator.Username))
		delete(m.spectators.users, spectator.UID)
		return
	}

	var p *PlayerReference
	var o *PlayerReference

	// assign the above variables, player and opponent of the closing socket
	if m.Player1 != nil && m.Player1.Socket == s {
		p = m.Player1

		if m.Player2 != nil && m.Player2.Socket != nil {
			o = m.Player2
		}

	} else if m.Player2 != nil && m.Player2.Socket == s {
		p = m.Player2

		if m.Player1 != nil && m.Player1.Socket != nil {
			o = m.Player1
		}
	}

	if p == nil {
		return
	}

	if o != nil {
		// let the opponent know that this player has disconnected
		o.Socket.Send(server.Message{
			Header: "opponent_disconnected",
		})
	}

	p.Socket = nil

	// if both players have disconnected, close match
	if (p == nil || p.Socket == nil) && (o == nil || o.Socket == nil) {
		m.quit <- true
	}

}
