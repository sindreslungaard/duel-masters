package match

import (
	"duel-masters/game/cnd"
	"duel-masters/server"
	"testing"
	"time"
)

const spawnCardTestUID = "spawn-card-test"

func TestSpawnCardReturnsCardAndErrors(t *testing.T) {
	registerSpawnCardTestConstructor()
	player := NewPlayer(&Match{}, 1)

	card, err := player.SpawnCard(spawnCardTestUID, HAND)
	if err != nil {
		t.Fatalf("expected SpawnCard to succeed: %v", err)
	}
	if card == nil || card.Zone != HAND {
		t.Fatalf("expected spawned card in hand, got %+v", card)
	}

	if _, err := player.SpawnCard("missing-spawn-card-constructor", HAND); err == nil {
		t.Fatal("expected an unknown card constructor to return an error")
	}
	if _, err := player.SpawnCard(spawnCardTestUID, BATTLEZONE); err == nil {
		t.Fatal("expected an unsupported spawn zone to return an error")
	}
}

func TestSpawnCardsToGivenZonesInitializesCardConditions(t *testing.T) {
	registerSpawnCardTestConstructor()
	m := newSpawnCardTestMatch()
	defer m.Dispose()

	spawnCardsToGivenZones(m, m.Player1.Player, []string{spawnCardTestUID}, []string{HAND}, []string{"/add", spawnCardTestUID})

	hand, err := m.Player1.Player.Container(HAND)
	if err != nil {
		t.Fatalf("expected hand container: %v", err)
	}
	if len(hand) != 1 {
		t.Fatalf("expected one spawned card, got %d", len(hand))
	}
	if !hand[0].HasCondition(cnd.Creature) {
		t.Fatal("expected the spawned card's UntapStep conditions to be initialized")
	}
}

func registerSpawnCardTestConstructor() {
	AddCard(spawnCardTestUID, func(card *Card) {
		card.Name = "Spawn Card Test"
		card.Use(func(card *Card, ctx *Context) {
			if _, ok := ctx.Event.(*UntapStep); ok {
				card.AddCondition(cnd.Creature, nil, card.ID)
			}
		})
	})
}

func newSpawnCardTestMatch() *Match {
	system := NewSystem()
	m := system.NewMatch("spawn-card-test", "host", "Player 1", nil, "", "Player 2", nil, true, true, FormatDescriptor{}, "host")

	player1 := NewPlayer(m, 1)
	player2 := NewPlayer(m, 2)
	m.Player1 = NewPlayerReference(player1, server.NewSocket(&spawnCardTestConnection{}, m, "player-1", "Player 1"))
	m.Player2 = NewPlayerReference(player2, server.NewSocket(&spawnCardTestConnection{}, m, "player-2", "Player 2"))
	m.Turn = player1.Turn

	return m
}

type spawnCardTestConnection struct{}

func (*spawnCardTestConnection) SetReadLimit(int64)                {}
func (*spawnCardTestConnection) SetReadDeadline(time.Time) error   { return nil }
func (*spawnCardTestConnection) SetPongHandler(func(string) error) {}
func (*spawnCardTestConnection) ReadMessage() (int, []byte, error) { return 0, nil, nil }
func (*spawnCardTestConnection) SetWriteDeadline(time.Time) error  { return nil }
func (*spawnCardTestConnection) WriteMessage(int, []byte) error    { return nil }
func (*spawnCardTestConnection) WriteJSON(interface{}) error       { return nil }
func (*spawnCardTestConnection) Close() error                      { return nil }
