package test

import "duel-masters/game/match"

type TestScenario struct {
	match *match.Match
}

type TestScenarioOptions struct{}

func NewTestScenario(options TestScenarioOptions) *TestScenario {
	matchSystem := match.NewSystem(func(msg interface{}) {})
	m := matchSystem.NewMatch("test-scenario", "test-host", true, true, match.RegularFormat)

	p1 := match.NewPlayer(m, 1)
	m.Player1 = match.NewPlayerReference(p1, s)

	m.Start()

	return &TestScenario{
		match: m,
	}
}
