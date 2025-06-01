package tests

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"
)

func TestMatchInitialization(t *testing.T) {
	scenario := scenario.New(scenario.Options{})

	p1manazone, err := scenario.Match.Player1.Player.Container(match.MANAZONE)

	if err != nil {
		t.Error(err)
	}

	if len(p1manazone) <= 0 {
		t.Errorf("Expected p1 manazone to have cards in it")
	}
}
