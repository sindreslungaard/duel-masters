package scenario

import (
	"duel-masters/game/cards"
	"duel-masters/game/match"
	"duel-masters/server"
)

type TestScenario struct {
	Match *match.Match
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
	m.Player1 = match.NewPlayerReference(p1, server.NewSocket(NewMockConnection(), m))

	p2 := match.NewPlayer(m, 2)
	m.Player2 = match.NewPlayerReference(p2, server.NewSocket(NewMockConnection(), m))

	deck := []string{}
	for range 40 {
		deck = append(deck, "af3bc221-1cc2-4f58-83ea-2673ac2c66c5") // Immortal Baron, Vorg
	}

	p1.CreateDeck(deck)
	p2.CreateDeck(deck)

	p1.Ready = true
	p2.Ready = true

	m.Start()

	return &TestScenario{
		Match: m,
	}
}
