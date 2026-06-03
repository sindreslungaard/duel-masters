package scenario

import (
	"duel-masters/game/cards"
	"duel-masters/game/match"
	"duel-masters/server"
)

type TestScenario struct {
	Match *match.Match
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

	p1 := match.NewPlayer(m, 1)
	m.Player1 = match.NewPlayerReference(p1, server.NewSocket(NewMockConnection(), m, "1", "Player1"))

	p2 := match.NewPlayer(m, 2)
	m.Player2 = match.NewPlayerReference(p2, server.NewSocket(NewMockConnection(), m, "2", "Player2"))

	p1.CreateDeck(cloneDeck(config.deck))
	p2.CreateDeck(cloneDeck(config.deck))

	p1.Ready = true
	p2.Ready = true

	m.Start()

	return &TestScenario{
		Match: m,
	}
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
