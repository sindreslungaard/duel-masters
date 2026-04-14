package tests

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchInitialization(t *testing.T) {
	scenario := scenario.New(scenario.Options{})

	t.Run("shieldzone is initialized with cards", func(t *testing.T) {
		p1shields, err := scenario.Match.Player1.Player.Container(match.SHIELDZONE)
		if err != nil {
			t.Error(err)
		}

		p2shields, err := scenario.Match.Player2.Player.Container(match.SHIELDZONE)
		if err != nil {
			t.Error(err)
		}

		assert.NotEmpty(t, p1shields)
		assert.NotEmpty(t, p2shields)
	})

	t.Run("manazone is empty", func(t *testing.T) {
		p1mana, err := scenario.Match.Player1.Player.Container(match.MANAZONE)
		if err != nil {
			t.Error(err)
		}

		p2mana, err := scenario.Match.Player2.Player.Container(match.MANAZONE)
		if err != nil {
			t.Error(err)
		}

		assert.Empty(t, p1mana)
		assert.Empty(t, p2mana)
	})

}
