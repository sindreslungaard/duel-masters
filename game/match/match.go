package match

import (
	"context"
	"duel-masters/db"
	"duel-masters/game/cnd"
	"duel-masters/internal"
	"duel-masters/server"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// Match struct
type Match struct {
	ID                string                                   `json:"id"`
	MatchName         string                                   `json:"name"`
	HostID            string                                   `json:"-"`
	Player1           *PlayerReference                         `json:"-"`
	Player2           *PlayerReference                         `json:"-"`
	spectators        internal.ConcurrentDictionary[Spectator] `json:"-"`
	persistentEffects map[int]PersistentEffect
	Turn              byte `json:"-"`
	Started           bool `json:"started"`
	Visible           bool `json:"visible"`
	Step              interface{}

	PlayerSelectingToss string
	TossPrediction      int
	TossOutcome         int

	Matchmaking bool
	created     int64
	ending      bool
	closed      bool
	isFirstTurn bool
	startedAt   int64

	eventloop *EventLoop
	system    *MatchSystem
}

// Name just returns "match", obligatory for a hub
func (m *Match) Name() string {
	return "match"
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

	m.eventloop.stop()

	for _, spectator := range m.spectators.Iter() {
		if spectator.Socket == nil {
			continue
		}
		spectator.Socket.Close()
		spectator.Socket = nil
	}

	if m.Player1 != nil {
		m.Player1.Dispose()
	}

	if m.Player2 != nil {
		m.Player2.Dispose()
	}

	m.system.Matches.Remove(m.ID)

	logrus.Debugf("Closed match with id %s", m.ID)

	m.system.UpdateMatchList()

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
		CardID:        card.ID,
		FromShield:    fromShield,
		MatchPlayerID: m.getPlayerMatchId(card),
	}))

	m.BroadcastState()

}

func (m *Match) getPlayerMatchId(card *Card) byte {
	if card.Player == m.Player1.Player {
		return 1
	}
	return 2
}

// Battle handles a battle between two creatures
func (m *Match) Battle(attacker *Card, defender *Card, blocked bool) {

	attackerPower := m.GetPower(attacker, true)
	defenderPower := m.GetPower(defender, false)

	m.HandleFx(NewContext(m, &AttackConfirmed{CardID: attacker.ID, Player: false, Creature: true}))
	m.HandleFx(NewContext(m, &Battle{Attacker: attacker, AttackerPower: attackerPower, Defender: defender, DefenderPower: defenderPower, Blocked: blocked}))

	m.BroadcastState()

}

// Destroy sends the given card to its players graveyard
func (m *Match) Destroy(card *Card, source *Card, context CreatureDestroyedContext) {

	m.HandleFx(NewContext(m, &CreatureDestroyed{Card: card, Source: source, Context: context}))
	m.Chat("Server", fmt.Sprintf("%s (%v) was destroyed by %s", card.Name, m.GetPower(card, false), source.Name))

}

// MoveCard moves a card and sends a chat message about what source moved it
func (m *Match) MoveCard(card *Card, destination string, source *Card) {

	_, err := card.Player.MoveCard(card.ID, card.Zone, destination, source.ID)

	if err != nil {
		return
	}

	m.Chat("Server", fmt.Sprintf("%s was moved to %s %s by %s", m.FormatDisplayableCard(card), card.Player.Username(), destination, m.FormatDisplayableCard(source)))

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
func (m *Match) BreakShields(shields []*Card, source string) {

	if len(shields) < 1 {
		return
	}

	m.Chat("Server", fmt.Sprintf("%v of %v's shields were broken", len(shields), m.PlayerRef(shields[0].Player).Socket.User.Username))

	for _, shield := range shields {

		card, err := shield.Player.MoveCard(shield.ID, SHIELDZONE, HAND, source)

		if err != nil {
			continue
		}

		m.HandleFx(NewContext(m, &BrokenShieldEvent{CardID: card.ID, Source: source}))

		// Handle shield triggers
		if card.HasCondition(cnd.ShieldTrigger) {

			ctx := NewContext(m, &ShieldTriggerEvent{
				Card:   card,
				Source: source,
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

				m.HandleFx(NewContext(m, &ShieldTriggerPlayedEvent{
					Card:   card,
					Source: source,
				}))

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

	m.ending = true

	if m.Started {
		m.Broadcast(server.WarningMessage{
			Header:  "error",
			Message: winnerStr,
		})

		m.SaveMatchHistory(winner, false)
	}

	m.Dispose()

}

func (m *Match) SaveMatchHistory(winner *Player, wonByDisconnect bool) {
	// don't save if the match lasted less than a minute
	if m.startedAt > time.Now().Unix()-60 {
		return
	}

	p1id := ""
	p1deck := ""
	p2id := ""
	p2deck := ""

	if m.Player1 != nil {
		p1id = m.Player1.UID
		p1deck = m.Player1.DeckStr
	}

	if m.Player2 != nil {
		p2id = m.Player2.UID
		p2deck = m.Player2.DeckStr
	}

	duel := db.Duel{
		UID:             m.ID,
		Host:            p1id,
		HostDeck:        p1deck,
		Guest:           p2id,
		GuestDeck:       p2deck,
		Started:         m.startedAt,
		Ended:           time.Now().Unix(),
		WonByDisconnect: wonByDisconnect,
	}

	if winner != nil && m.Player1 != nil && m.Player1.Player == winner {
		duel.Winner = m.Player1.UID
	} else if winner != nil && m.Player2 != nil && m.Player2.Player == winner {
		duel.Winner = m.Player2.UID
	}

	_, err := db.Duels.InsertOne(context.Background(), duel)

	if err != nil {
		logrus.Error("Failed to save duel result to db", err)
	}
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

func (m *Match) FormatDisplayableCard(card *Card) string {
	return fmt.Sprintf("(%s;%s)", card.Name, card.ImageID)
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
	for _, spectator := range m.spectators.Iter() {
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

	for _, spectator := range m.spectators.Iter() {
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

// DefaultActionWarning sends an action warning with a predefined message
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

	// End the game  if any of the players are out of cards in their deck
	for _, p := range players {
		if len(p.Player.deck) < 1 {
			m.End(m.Opponent(p.Player), fmt.Sprintf("%s won by deck out!", m.Opponent(p.Player).Username()))
		}
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

	player.ActionState = PlayerActionState{
		resolved: false,
		data:     msg,
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

	player.ActionState = PlayerActionState{
		resolved: false,
		data:     msg,
	}

	m.PlayerRef(player).Socket.Send(msg)

}

// NewMultipartAction prompts the user to make a selection of the specified {string: []Cards}
func (m *Match) NewMultipartAction(player *Player, cards map[string][]*Card, minSelections int, maxSelections int, text string, cancellable bool) {
	m.newMultipartActionBase(player, cards, minSelections, maxSelections, text, cancellable, false)
}

// NewMultipartAction prompts the user to make a selection of the specified {string: []Cards}
func (m *Match) NewMultipartActionBackside(player *Player, cards map[string][]*Card, minSelections int, maxSelections int, text string, cancellable bool) {
	m.newMultipartActionBase(player, cards, minSelections, maxSelections, text, cancellable, true)
}

func (m *Match) newMultipartActionBase(player *Player, cards map[string][]*Card, minSelections int, maxSelections int, text string, cancellable bool, backsideOnly bool) {

	cardMap := make(map[string][]server.CardState)

	for key, cards := range cards {
		cardMap[key] = denormalizeCards(cards, backsideOnly)
	}

	msg := &server.MultipartActionMessage{
		Header:        "action",
		Cards:         cardMap,
		Text:          text,
		MinSelections: minSelections,
		MaxSelections: maxSelections,
		Cancellable:   cancellable,
	}

	player.ActionState = PlayerActionState{
		resolved: false,
		data:     msg,
	}

	m.PlayerRef(player).Socket.Send(msg)

}

// CloseAction closes the card selection popup for the given player
func (m *Match) CloseAction(p *Player) {
	p.ActionState.resolved = true
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

// Coin toss to decide who starts the game
func (m *Match) CoinToss() {
	var PlayerToChooseToss *PlayerReference
	var PlayerToAwaitToss *PlayerReference

	if m.PlayerSelectingToss == "" {
		if rand.Intn(100) >= 50 {
			m.PlayerSelectingToss = m.Player2.UID

			PlayerToChooseToss = m.Player2
			PlayerToAwaitToss = m.Player1
		} else {
			m.PlayerSelectingToss = m.Player1.UID

			PlayerToChooseToss = m.Player1
			PlayerToAwaitToss = m.Player2
		}
	}

	m.Chat("Server", fmt.Sprintf("Coin toss to be selected by %s", PlayerToChooseToss.Username))

	PlayerToAwaitToss.Socket.Send(server.Message{
		Header: "toss_being_chosen",
	})

	PlayerToChooseToss.Socket.Send(server.Message{
		Header: "choose_toss",
	})
}

// Start starts the match
func (m *Match) Start() {

	m.Started = true
	m.startedAt = time.Now().Unix()

	m.system.UpdateMatchList()

	m.Player1.Player.ShuffleDeck()
	m.Player2.Player.ShuffleDeck()

	m.Player1.Player.InitShieldzone()
	m.Player2.Player.InitShieldzone()

	m.Player1.Player.DrawCards(5)
	m.Player2.Player.DrawCards(5)

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

	if card, err := p.Player.MoveCard(cardID, HAND, MANAZONE, cardID); err == nil {
		p.Player.HasChargedMana = true
		m.BroadcastState()
		m.Chat("Server", fmt.Sprintf("%s was added to %s's manazone",
			m.FormatDisplayableCard(card),
			p.Socket.User.Username))
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

// AttackCreature is called when the player attempts to attack an opposing creature
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

func (m *Match) TapAbility(p *PlayerReference, cardID string) {
	_, err := p.Player.GetCard(cardID, BATTLEZONE)

	if err != nil {
		Warn(p, "The creature you tried to use a tap ability for is not in the battlezone")
		return
	}

	ctx := NewContext(m, &TapAbility{
		CardID: cardID,
	})

	m.HandleFx(ctx)

	if !ctx.Cancelled() {
		// Tap abilities can only be used during attack step
		// https://duelmasters.fandom.com/wiki/Step#Step_7_(Attack_step)
		if _, ok := m.Step.(*AttackStep); !ok {
			m.Step = &AttackStep{}
		}

		p.Player.CanChargeMana = false
	}

	m.BroadcastState()
}

// Parse handles websocket messages in this Hub
func (m *Match) Parse(s *server.Socket, data []byte) {

	defer internal.Recover()

	var message server.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return
	}

	switch message.Header {

	case "mpong":
		m.eventloop.schedule(func() {
			p, err := m.PlayerForSocket(s)

			if err != nil {
				return
			}

			p.LastPong = time.Now().Unix()
		}, ParallelEvent)

	case "join_match":
		m.eventloop.schedule(func() {
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

					// if currently pending an action from this player
					if !m.Player1.Player.ActionState.resolved {
						s.Send(m.Player1.Player.ActionState.data)
					}

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

					// if currently pending an action from this player
					if !m.Player2.Player.ActionState.resolved {
						s.Send(m.Player2.Player.ActionState.data)
					}

					m.Chat("Server", s.User.Username+" reconnected")

					return

				} else {
					// Spectators

					spectator, ok := m.spectators.Find(s.User.UID)

					// this user is already spectating, swap connection to new one
					if ok {
						spectator.Socket.Send(server.WarningMessage{
							Header:  "error",
							Message: "You started spectating from a new connection, closing this one...",
						})
						// this removes the existing spectator from the match
						spectator.Socket.Close()
					}

					m.spectators.Add(s.User.UID, &Spectator{
						UID:      s.User.UID,
						Username: s.User.Username,
						Color:    s.User.Color,
						Socket:   s,
						LastPong: time.Now().Unix(),
					})

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

				cur, err := db.Decks.Find(context.TODO(), bson.M{
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

				m.Chat("Server", fmt.Sprintf("%s started the game", m.Player1.Username))
				m.Chat("Server", fmt.Sprintf("%s joined the game", m.Player2.Username))

				player1decks := make([]db.LegacyDeck, 0)
				player2decks := make([]db.LegacyDeck, 0)

				for cur.Next(context.TODO()) {

					var deck db.Deck

					if err := cur.Decode(&deck); err != nil {
						continue
					}

					legacyDeck, err := ConvertToLegacyDeck(deck)
					if err != nil {
						continue
					}

					if deck.Owner == m.Player1.Socket.User.UID || deck.Standard {
						player1decks = append(player1decks, legacyDeck)
					}

					if deck.Owner == m.Player2.Socket.User.UID || deck.Standard {
						player2decks = append(player2decks, legacyDeck)
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

			m.system.UpdateMatchList()

		}, ParallelEvent)

	case "chat":
		m.eventloop.schedule(func() {

			// Allow other sockets than player1 and player2 to chat?

			var msg struct {
				Message string `json:"message"`
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			m.handleAdminMesseges(msg.Message, s.User)

			if s.User.Chatblocked {
				s.Send(&server.ChatMessage{
					Header:  "chat",
					Message: msg.Message,
					Sender:  s.User.Username,
					Color:   s.User.Color,
				})
			} else {
				m.ColorChat(s.User.Username, msg.Message, s.User.Color)
			}
		}, ParallelEvent)

	case "toss_chosen":
		{
			m.eventloop.schedule(func() {
				if m.Player1.Player.Ready && m.Player2.Player.Ready {
					if !m.Started && m.PlayerSelectingToss != "" && m.PlayerSelectingToss == s.User.UID {

						var msg struct {
							Prediction int `json:"prediction"`
						}

						if err := json.Unmarshal(data, &msg); err != nil {
							return
						}

						if msg.Prediction == 1 || msg.Prediction == -1 {
							if msg.Prediction == 1 {
								m.Chat("Server", fmt.Sprintf("%s has called heads", s.User.Username))
							} else if msg.Prediction == -1 {
								m.Chat("Server", fmt.Sprintf("%s has called tails", s.User.Username))
							}

							var outcome = rand.Intn(100)
							var heads = (outcome >= 50)
							var wasTossDecisionCorrect bool

							m.TossPrediction = msg.Prediction

							if heads {
								m.TossOutcome = 1
							} else {
								m.TossOutcome = -1
							}

							if (heads && msg.Prediction == 1) || (!heads && msg.Prediction == -1) {
								wasTossDecisionCorrect = true

								m.Chat("Server", fmt.Sprintf("%s won the coin toss", s.User.Username))
							} else {
								wasTossDecisionCorrect = false

								m.Chat("Server", fmt.Sprintf("%s lost the coin toss", s.User.Username))
							}

							var PlayerToChooseToss *PlayerReference
							var PlayerToAwaitToss *PlayerReference

							var PlayerToChooseTurnSelection *PlayerReference
							var PlayerToAwaitTurnSelection *PlayerReference

							if m.Player2.UID == m.PlayerSelectingToss {
								PlayerToChooseToss = m.Player2
								PlayerToAwaitToss = m.Player1
							} else {
								PlayerToChooseToss = m.Player1
								PlayerToAwaitToss = m.Player2
							}

							if wasTossDecisionCorrect {
								PlayerToChooseTurnSelection = PlayerToChooseToss
								PlayerToAwaitTurnSelection = PlayerToAwaitToss
							} else {
								PlayerToChooseTurnSelection = PlayerToAwaitToss
								PlayerToAwaitTurnSelection = PlayerToChooseToss
							}

							m.Chat("Server", fmt.Sprintf("%s to pick who goes first", PlayerToChooseTurnSelection.Username))

							PlayerToAwaitTurnSelection.Socket.Send(server.Message{
								Header: "turn_being_chosen",
							})

							PlayerToChooseTurnSelection.Socket.Send(server.Message{
								Header: "choose_turn",
							})
						}
					}
				}

				return
			}, SequentialEvent)
		}

	case "turn_chosen":
		{
			m.eventloop.schedule(func() {
				if m.Player1.Player.Ready && m.Player2.Player.Ready {
					if !m.Started && m.PlayerSelectingToss != "" && m.TossOutcome != 0 && m.TossPrediction != 0 {

						var PlayerToChooseToss *PlayerReference
						var PlayerToAwaitToss *PlayerReference

						var PlayerWhoWonToss *PlayerReference

						if m.Player2.UID == m.PlayerSelectingToss {
							PlayerToChooseToss = m.Player2
							PlayerToAwaitToss = m.Player1
						} else {
							PlayerToChooseToss = m.Player1
							PlayerToAwaitToss = m.Player2
						}

						if m.TossPrediction == m.TossOutcome {
							PlayerWhoWonToss = PlayerToChooseToss
						} else {
							PlayerWhoWonToss = PlayerToAwaitToss
						}

						if PlayerWhoWonToss.UID == s.User.UID {
							var msg struct {
								Player int `json:"player"`
							}

							if err := json.Unmarshal(data, &msg); err != nil {
								return
							}

							var isSelectingUserPlayer1 bool

							if m.Player1.UID == s.User.UID {
								isSelectingUserPlayer1 = true
							}

							// Socket user goes first
							if msg.Player == 1 {
								m.Chat("Server", fmt.Sprintf("%s decided to play first", PlayerWhoWonToss.Username))

								// match.turn is initialized as 1, so we only need to change it to 2
								// The opposite of what's defined here will start because BeginNewTurn() changes it
								if isSelectingUserPlayer1 {
									m.Turn = 2
								}
							} else if msg.Player == -1 {
								m.Chat("Server", fmt.Sprintf("%s decided to play second", PlayerWhoWonToss.Username))

								// match.turn is initialized as 1, so we only need to change it to 2
								// The opposite of what's defined here will start because BeginNewTurn() changes it
								if !isSelectingUserPlayer1 {
									m.Turn = 2
								}
							}

							m.Start()
						}
					}
				}

				return

			}, SequentialEvent)
		}

	case "choose_deck":
		m.eventloop.schedule(func() {

			if m.Started {
				return
			}

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

			if err := db.Decks.FindOne(context.TODO(), bson.M{"uid": msg.UID}).Decode(&deck); err != nil {
				return
			}

			p.DeckStr = deck.Cards

			legacyDeck, err := ConvertToLegacyDeck(deck)

			if err != nil {
				return
			}

			p.Player.CreateDeck(legacyDeck.Cards)

			m.Chat("Server", fmt.Sprintf("%s has chosen their deck", s.User.Username))

			p.Player.Ready = true

			if m.Player1.Player.Ready && m.Player2.Player.Ready {
				m.CoinToss()
			}

		}, SequentialEvent)

	case "add_to_manazone":
		m.eventloop.schedule(func() {

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

		}, SequentialEvent)

	case "end_turn":
		m.eventloop.schedule(func() {

			p, err := m.PlayerForSocket(s)

			if err != nil {
				return
			}

			if m.Turn != p.Player.Turn {
				return
			}

			m.EndTurn()

		}, SequentialEvent)

	case "add_to_playzone":
		m.eventloop.schedule(func() {

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

		}, SequentialEvent)

	case "action":
		m.eventloop.schedule(func() {

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

			// Drain the player action channel to prevent malicious exhaustion of goroutines
			// as the msg is sent to the channel in a new grouroutine for every event
		Drain:
			for {
				select {
				case <-p.Player.Action:
				default:
					break Drain
				}
			}

			p.Player.Action <- msg

		}, ParallelEvent)

	case "attack_player":
		m.eventloop.schedule(func() {

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

		}, SequentialEvent)

	case "attack_creature":
		m.eventloop.schedule(func() {

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

		}, SequentialEvent)

	case "tap_ability":
		m.eventloop.schedule(func() {

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

			m.TapAbility(p, msg.ID)

		}, SequentialEvent)

	default:
		logrus.Debugf("Received message in incorrect format: %v", string(data))

	}

}

// OnSocketClose is called when a socket disconnects
func (m *Match) OnSocketClose(s *server.Socket) {

	if m.closed {
		return
	}

	// is this a spectator leaving?
	spectator, ok := m.spectators.Find(s.User.UID)

	if ok {
		m.Chat("Server", fmt.Sprintf("%s stopped spectating", spectator.Username))
		m.spectators.Remove(spectator.UID)
		return
	}

	var p *PlayerReference // this is the disconnecting user
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

	if !m.Started {
		logrus.Debug(fmt.Sprintf("Disconnected %s Host %s", s.User.UID, m.HostID))

		// If the host left before the game started, close the game
		if p.UID == m.HostID {
			logrus.Debug("Host socket left before match started, closing match.")

			if o != nil {
				m.Chat("Server", fmt.Sprintf("%s disconnected from the game", p.Username))
			}

			m.Dispose()
			return
		}

		// This was the joinee who disconnected so let's reset things
		if m.Player1.Socket == s {
			logrus.Debug("Player1 socket left before match started, resetting Player1.")
			m.Player1 = nil
		}

		if m.Player2.Socket == s {
			logrus.Debug("Player2 socket left before match started, resetting Player2.")
			m.Player2 = nil
		}

		o.Player.DestroyDeck()
		o.Player.Ready = false

		m.PlayerSelectingToss = ""
		m.TossOutcome = 0
		m.TossPrediction = 0

		m.system.UpdateMatchList()
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

	// if both players have disconnected, close match
	if (p == nil || p.Socket.IsClosed()) && (o == nil || o.Socket.IsClosed()) {
		logrus.Debug("Both players left the match. Closing the match.")

		if p != nil {
			m.SaveMatchHistory(p.Player, true)
		}
		m.Dispose()
	}

}

func (p *PlayerReference) Dispose() {
	defer internal.Recover()

	if p.Player != nil {
		p.Player.Dispose()
	}

	if p.Socket != nil {
		p.Socket.Close()
	}
}

func (m *Match) handleAdminMesseges(message string, user db.User) {

	hasRights := false

	for _, permission := range user.Permissions {
		if permission == "admin" {
			hasRights = true
		}
	}

	if !hasRights {
		return
	}

	currentPlayer := m.CurrentPlayer().Player

	if string(message[0:4]) == "/add" {

		// Spawn card in hand
		currentPlayer.SpawnCard(string(message[5:]), HAND)
		m.BroadcastState()
		return
	}

	// /mana hand - will add all the cards in your hand to the manazone
	// /mana red 4 - add 4 fire mana
	// /mana n - add 1 nature mana
	// /mana imageID 3 - add 3 mana of that specific card
	if string(message[0:5]) == "/mana" {

		msgParts := strings.Split(message, " ")

		var manaToAdd string
		switch msgParts[1] {
		case "hand":
			for _, c := range currentPlayer.hand {
				currentPlayer.MoveCard(c.ID, HAND, MANAZONE, "cmd /mana")
				m.Chat("Server", fmt.Sprintf("%s was moved to %s's mana zone by an admin command", c.Name, user.Username))
			}
		case "fire", "f", "red":
			manaToAdd = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5"
		case "water", "w", "blue":
			manaToAdd = "9781089f-1aa9-4a75-b106-35e9d431e31d"
		case "light", "l", "yellow":
			manaToAdd = "7b58e8c2-0b1e-4ef5-812f-e667c2092c73"
		case "darkness", "d", "black":
			manaToAdd = "e2b992ee-91a3-49d3-8228-7be60a0b9ec5"
		case "nature", "n", "green":
			manaToAdd = "1d72eb3e-5185-449a-a16f-391bd2338343"
		default:
			manaToAdd = msgParts[1]
		}

		if manaToAdd != "" {
			count := 1
			if len(msgParts) > 2 {
				number, err := strconv.Atoi(msgParts[2])
				if err == nil && number <= 5 || number >= 1 {
					count = number
				}
			}

			for i := 0; i < count; i++ {
				currentPlayer.SpawnCard(manaToAdd, MANAZONE)
			}
		}

		m.BroadcastState()

		return
	}

	if string(message[0:7]) == "/shield" {

		// Spawn shield
		currentPlayer.SpawnCard(string(message[8:]), SHIELDZONE)
		m.BroadcastState()
		return
	}

	if string(message[0:5]) == "/deck" {

		// Spawn card in deck
		m.CurrentPlayer().Player.SpawnCard(string(message[6:]), DECK)
		m.BroadcastState()
		return
	}
}
