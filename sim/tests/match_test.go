package tests

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchInitialization(t *testing.T) {
	scn := scenario.New()

	t.Run("shieldzone is initialized with cards", func(t *testing.T) {
		p1shields, err := scn.Match.Player1.Player.Container(match.SHIELDZONE)
		if err != nil {
			t.Error(err)
		}

		p2shields, err := scn.Match.Player2.Player.Container(match.SHIELDZONE)
		if err != nil {
			t.Error(err)
		}

		assert.NotEmpty(t, p1shields)
		assert.NotEmpty(t, p2shields)
	})

	t.Run("manazone is empty", func(t *testing.T) {
		p1mana, err := scn.Match.Player1.Player.Container(match.MANAZONE)
		if err != nil {
			t.Error(err)
		}

		p2mana, err := scn.Match.Player2.Player.Container(match.MANAZONE)
		if err != nil {
			t.Error(err)
		}

		assert.Empty(t, p1mana)
		assert.Empty(t, p2mana)
	})

	t.Run("custom deck preserves requested card counts after setup", func(t *testing.T) {
		customScenario := scenario.New(
			scenario.WithDeck(
				scenario.DeckEntry{UID: "c5a869f4-a959-4667-a352-92df5369e0b9", Count: 5},
				scenario.DeckEntry{UID: "af3bc221-1cc2-4f58-83ea-2673ac2c66c5", Count: 35},
			),
		)

		cards := customScenario.Match.Player1.Player.Cards()
		counts := map[string]int{}
		for _, card := range cards {
			counts[card.ImageID]++
		}

		assert.Len(t, cards, 40)
		assert.Equal(t, 5, counts["c5a869f4-a959-4667-a352-92df5369e0b9"])
		assert.Equal(t, 35, counts["af3bc221-1cc2-4f58-83ea-2673ac2c66c5"])
	})

}
