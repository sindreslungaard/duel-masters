package scenario

import (
	"duel-masters/game/cards"
	"duel-masters/game/match"
	"duel-masters/server"
)

type TestScenario struct {
	Match   *match.Match
	Client1 *MockClient
	Client2 *MockClient
}

type Options struct{}

func New(options Options) *TestScenario {
	for _, set := range cards.Sets {
		for uid, ctor := range *set {
			if ctor == nil {
				continue
			}
			match.AddCard(uid, ctor)
		}
	}

	matchSystem := match.NewSystem(func(msg interface{}) {})
	m := matchSystem.NewMatch("test-scenario", "test-host", true, true, match.RegularFormat)

	p1 := match.NewPlayer(m, 1)
	s1 := server.NewSocket(NewMockConnection(), m)
	c1 := NewMockClient(s1, m)
	s1.SetReady(true)
	m.Player1 = match.NewPlayerReference(p1, s1)

	p2 := match.NewPlayer(m, 2)
	s2 := server.NewSocket(NewMockConnection(), m)
	c2 := NewMockClient(s2, m)
	s2.SetReady(true)
	m.Player2 = match.NewPlayerReference(p2, s2)

	deck := []string{}
	for range 40 {
		deck = append(deck, "af3bc221-1cc2-4f58-83ea-2673ac2c66c5") // Immortal Baron, Vorg
	}

	p1.CreateDeck(deck)
	p2.CreateDeck(deck)

	p1.Ready = true
	p2.Ready = true

	m.Turn = 2 // so that player1 starts
	m.Start()

	return &TestScenario{
		Match:   m,
		Client1: c1,
		Client2: c2,
	}
}
