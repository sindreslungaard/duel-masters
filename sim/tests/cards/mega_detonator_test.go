package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const megaDetonatorUID = "e6c76df1-24c8-4125-9f9f-8ae3b2bc61f6"

func TestMegaDetonator(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(megaDetonatorUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Mega Detonator", spell.Name)
		assert.Equal(t, 2, spell.ManaCost)
		assert.Equal(t, []string{civ.Fire}, spell.Civs)
	})

	t.Run("allows discarding every remaining card in hand", func(t *testing.T) {
		// Regression test: the maximum discard count used to be computed as
		// (hand size) - 1, to compensate for Mega Detonator counting itself
		// while still in hand. Now that fx.Spell moves it to the graveyard
		// before this effect runs, that subtraction undercounts the true
		// remaining hand by one.
		scn, player, _ := setupDuel(t)

		emptyHand(t, player, "mega_detonator_test_setup")
		fillers := make([]*match.Card, 0, 3)
		for range 3 {
			c, err := player.Player.SpawnCard(scowlingTomatoUID, match.HAND)
			require.NoError(t, err)
			fillers = append(fillers, c)
		}

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		castSpell(t, scn, player, megaDetonatorUID)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err, "expected the discard-selection prompt to be open")
		assert.Equal(t, 3, action.MaxSelections, "all 3 remaining hand cards should be discardable, not 2")

		fillerIDs := make([]string, 0, len(fillers))
		for _, filler := range fillers {
			fillerIDs = append(fillerIDs, filler.ID)
		}
		require.NoError(t, scn.SubmitAction(player, fillerIDs...))
		require.NoError(t, scn.WaitForEventLoop())

		for _, filler := range fillers {
			assert.Equal(t, match.GRAVEYARD, filler.Zone)
		}
	})
}
