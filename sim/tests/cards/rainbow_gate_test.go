package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	rainbowGateUID          = "3a071f56-5ad3-445d-b223-6f76685d843a"
	rainbowGateMulticolorID = "ef6d1314-1005-4b5d-9afe-827ffe3ba58a" // Aqua Skydiver (light/water)
	rainbowGateSetupSrc     = "rainbow_gate_test_setup"
)

func TestRainbowGate(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(rainbowGateUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Rainbow Gate", spell.Name)
		assert.Equal(t, 2, spell.ManaCost)
		assert.Equal(t, []string{civ.Nature}, spell.Civs)
	})

	t.Run("it takes a multicolored creature out of the deck", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// Seeded below the top so the draw step cannot take the very card the
		// search is supposed to find.
		player.Player.DestroyDeck()
		for range 3 {
			_, err := player.Player.SpawnCard(immortalBaronVorgUID, match.DECK)
			require.NoError(t, err)
		}
		multicolored, err := player.Player.SpawnCard(rainbowGateMulticolorID, match.DECK)
		require.NoError(t, err)
		for range 3 {
			_, err := player.Player.SpawnCard(immortalBaronVorgUID, match.DECK)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, rainbowGateUID)
		answerInTurn(t, scn, player, multicolored.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, multicolored.Zone)
	})

	t.Run("a single colour creature is not a legal choice", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		player.Player.DestroyDeck()
		for range 6 {
			_, err := player.Player.SpawnCard(immortalBaronVorgUID, match.DECK)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		deckBefore, err := player.Player.Container(match.DECK)
		require.NoError(t, err)

		castSpell(t, scn, player, rainbowGateUID)

		// The search shows the whole deck, so the prompt still opens; there is
		// simply nothing in it that can be taken.
		cancelInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		deck, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		assert.Len(t, deck, len(deckBefore), "nothing in the deck could be taken")
	})
}
