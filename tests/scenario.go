package test

import "duel-masters/game/match"

type TestScenario struct {
	match *match.Match
}

type TestScenarioOptions struct{}

func NewTestScenario(options TestScenarioOptions) *TestScenario {
	matchSystem := match.NewSystem(func(msg interface{}) {})
	match := matchSystem.NewMatch("test-scenario", "test-host", true, true, match.RegularFormat)

	match.Start()

	return &TestScenario{
		match,
	}
}
