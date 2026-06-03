package scenario

import (
	"duel-masters/game/cards"
	"duel-masters/game/match"
	"duel-masters/server"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const actionTimeout = 2 * time.Second

type TestScenario struct {
	Match       *match.Match
	connections map[*match.PlayerReference]*MockConnection
}

type Option func(*scenarioConfig)

type scenarioConfig struct {
	deck []string
}

type DeckEntry struct {
	UID   string
	Count int
}

func WithDeck(entries ...DeckEntry) Option {
	return func(opts *scenarioConfig) {
		deck := make([]string, 0)

		for _, entry := range entries {
			for range entry.Count {
				deck = append(deck, entry.UID)
			}
		}

		opts.deck = deck
	}
}

func New(options ...Option) *TestScenario {
	for _, set := range cards.Sets {
		for uid, ctor := range *set {
			if ctor == nil {
				continue
			}
			match.AddCard(uid, ctor)
		}
	}

	matchSystem := match.NewSystem(func(msg interface{}) {})
	m := matchSystem.NewMatch("test-scenario", "test-host", []string{}, "", []string{}, true, true, match.RegularFormat)

	config := scenarioConfig{deck: defaultDeck()}
	for _, option := range options {
		option(&config)
	}

	p1Conn := NewMockConnection()
	p1 := match.NewPlayer(m, 1)
	m.Player1 = match.NewPlayerReference(p1, server.NewSocket(p1Conn, m, "1", "Player1"))

	p2Conn := NewMockConnection()
	p2 := match.NewPlayer(m, 2)
	m.Player2 = match.NewPlayerReference(p2, server.NewSocket(p2Conn, m, "2", "Player2"))

	p1.CreateDeck(cloneDeck(config.deck))
	p2.CreateDeck(cloneDeck(config.deck))

	p1.Ready = true
	p2.Ready = true

	m.Start()

	return &TestScenario{
		Match: m,
		connections: map[*match.PlayerReference]*MockConnection{
			m.Player1: p1Conn,
			m.Player2: p2Conn,
		},
	}
}

func (s *TestScenario) FindCard(player *match.Player, zone string, imageID string) (*match.Card, error) {
	cards, err := player.Container(zone)
	if err != nil {
		return nil, err
	}

	for _, card := range cards {
		if card.ImageID == imageID {
			return card, nil
		}
	}

	return nil, fmt.Errorf("card %s not found in %s", imageID, zone)
}

func (s *TestScenario) ActionPlayCard(player *match.PlayerReference, cardID string) error {
	conn, err := s.connectionFor(player)
	if err != nil {
		return err
	}

	card, err := player.Player.GetCard(cardID, match.HAND)
	if err != nil {
		return err
	}

	messageCount := conn.JSONWriteCount()

	if err := s.send(player, struct {
		Header string `json:"header"`
		ID     string `json:"virtualId"`
	}{
		Header: "add_to_playzone",
		ID:     cardID,
	}); err != nil {
		return err
	}

	action, err := s.waitForActionMessage(player, messageCount)
	if err != nil {
		return err
	}

	selection, err := s.selectPlayableMana(card, action)
	if err != nil {
		return err
	}

	if err := s.SubmitAction(player, selection...); err != nil {
		return err
	}

	return s.waitFor(func() bool { return card.Zone != match.HAND })
}

func (s *TestScenario) ActionEndTurn(player *match.PlayerReference) error {
	conn, err := s.connectionFor(player)
	if err != nil {
		return err
	}

	turn := s.Match.Turn
	messageCount := conn.JSONWriteCount()

	if err := s.send(player, struct {
		Header string `json:"header"`
	}{
		Header: "end_turn",
	}); err != nil {
		return err
	}

	return s.waitFor(func() bool {
		return s.Match.Turn != turn || conn.JSONWriteCount() > messageCount
	})
}

func (s *TestScenario) SubmitAction(player *match.PlayerReference, cardIDs ...string) error {
	return s.send(player, struct {
		Header string   `json:"header"`
		Cards  []string `json:"cards"`
		Count  int      `json:"count"`
		Cancel bool     `json:"cancel"`
	}{
		Header: "action",
		Cards:  cardIDs,
		Count:  len(cardIDs),
		Cancel: false,
	})
}

// CancelAction answers "no" or cancels a pending action prompt (e.g. a binary yes/no question).
func (s *TestScenario) CancelAction(player *match.PlayerReference) error {
	return s.send(player, struct {
		Header string   `json:"header"`
		Cards  []string `json:"cards"`
		Count  int      `json:"count"`
		Cancel bool     `json:"cancel"`
	}{
		Header: "action",
		Cards:  []string{},
		Count:  0,
		Cancel: true,
	})
}

// MessageHeaders returns the header field of every JSON message sent to the player so far,
// starting from position since. Useful for debugging test failures.
func (s *TestScenario) MessageHeaders(player *match.PlayerReference, since int) ([]string, error) {
	conn, err := s.connectionFor(player)
	if err != nil {
		return nil, err
	}
	var headers []string
	for _, raw := range conn.JSONMessagesSince(since) {
		var h server.Message
		if json.Unmarshal([]byte(raw), &h) == nil {
			headers = append(headers, h.Header)
		}
	}
	return headers, nil
}

// MessageCount returns the number of JSON messages the server has sent to the player so far.
// Capture this value BEFORE an action so that WaitForMessage can find the response.
func (s *TestScenario) MessageCount(player *match.PlayerReference) (int, error) {
	conn, err := s.connectionFor(player)
	if err != nil {
		return 0, err
	}
	return conn.JSONWriteCount(), nil
}

// WaitForMessage blocks until a server message whose header matches one of the supplied values
// appears at a position >= since in the player's received message log.
// Always capture since = MessageCount(player) BEFORE the action that triggers the message.
func (s *TestScenario) WaitForMessage(player *match.PlayerReference, since int, headers ...string) error {
	conn, err := s.connectionFor(player)
	if err != nil {
		return err
	}

	return s.waitFor(func() bool {
		for _, raw := range conn.JSONMessagesSince(since) {
			var h server.Message
			if err := json.Unmarshal([]byte(raw), &h); err != nil {
				continue
			}
			for _, header := range headers {
				if h.Header == header {
					return true
				}
			}
		}
		return false
	})
}

func (s *TestScenario) connectionFor(player *match.PlayerReference) (*MockConnection, error) {
	conn, ok := s.connections[player]
	if !ok {
		return nil, errors.New("mock connection not found for player")
	}

	return conn, nil
}

func (s *TestScenario) send(player *match.PlayerReference, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	s.Match.Parse(player.Socket, data)
	return nil
}

func (s *TestScenario) waitForActionMessage(player *match.PlayerReference, start int) (*server.ActionMessage, error) {
	conn, err := s.connectionFor(player)
	if err != nil {
		return nil, err
	}

	var action *server.ActionMessage

	err = s.waitFor(func() bool {
		for _, raw := range conn.JSONMessagesSince(start) {
			var header server.Message
			if err := json.Unmarshal([]byte(raw), &header); err != nil {
				continue
			}

			if header.Header != "action" {
				continue
			}

			candidate := &server.ActionMessage{}
			if err := json.Unmarshal([]byte(raw), candidate); err != nil {
				continue
			}

			action = candidate
			return true
		}

		return false
	})
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (s *TestScenario) selectPlayableMana(card *match.Card, action *server.ActionMessage) ([]string, error) {
	blocked := make(map[string]struct{}, len(action.UnselectableCards))
	for _, unselectable := range action.UnselectableCards {
		blocked[unselectable.CardID] = struct{}{}
	}

	available := make([]*match.Card, 0, len(action.Cards))
	for _, candidate := range action.Cards {
		if _, ok := blocked[candidate.CardID]; ok {
			continue
		}

		mana, err := card.Player.GetCard(candidate.CardID, match.MANAZONE)
		if err != nil {
			continue
		}

		available = append(available, mana)
	}

	for selections := action.MinSelections; selections <= action.MaxSelections; selections++ {
		if picked, ok := pickPlayableMana(card, available, selections, 0, nil); ok {
			return picked, nil
		}
	}

	return nil, fmt.Errorf("could not find playable mana selection for %s", card.Name)
}

func (s *TestScenario) waitFor(predicate func() bool) error {
	deadline := time.Now().Add(actionTimeout)
	for time.Now().Before(deadline) {
		if predicate() {
			return nil
		}

		time.Sleep(time.Millisecond)
	}

	return fmt.Errorf("timed out waiting for scenario action after %s", actionTimeout)
}

func pickPlayableMana(card *match.Card, available []*match.Card, target int, start int, current []*match.Card) ([]string, bool) {
	if len(current) == target {
		if !card.Player.CanPlayCard(card, current) {
			return nil, false
		}

		selection := make([]string, len(current))
		for i, mana := range current {
			selection[i] = mana.ID
		}

		return selection, true
	}

	remaining := target - len(current)
	for i := start; i <= len(available)-remaining; i++ {
		if selection, ok := pickPlayableMana(card, available, target, i+1, append(current, available[i])); ok {
			return selection, true
		}
	}

	return nil, false
}

func defaultDeck() []string {
	deck := make([]string, 0, 40)
	for range 40 {
		deck = append(deck, "af3bc221-1cc2-4f58-83ea-2673ac2c66c5") // Immortal Baron, Vorg
	}

	return deck
}

func cloneDeck(deck []string) []string {
	cloned := make([]string, len(deck))
	copy(cloned, deck)
	return cloned
}
