package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	uncannyTurnipUID       = "a6bd38a4-734f-4242-bbd2-93942396d1cb"
	uncannyTurnipFillerUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	uncannyTurnipSetupSrc  = "uncanny_turnip_test_setup"
)

func TestUncannyTurnip(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		// Kept in hand: it prompts as it arrives once its ability is switched on.
		card := spawnForLater(t, player, uncannyTurnipUID)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Uncanny Turnip", 1000, 2, []string{civ.Nature})
		assert.True(t, card.HasFamily(family.WildVeggies))
		assert.True(t, card.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count it feeds mana and takes a creature back", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)

		inMana, err := player.Player.SpawnCard(uncannyTurnipFillerUID, match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		topBefore := player.Player.PeekDeck(1)
		require.Len(t, topBefore, 1)

		summonWithOwnMana(t, scn, player, uncannyTurnipUID)

		// Two creatures in mana now, so the return has to be chosen.
		answerInTurn(t, scn, player, inMana.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.MANAZONE, topBefore[0].Zone, "the top of the deck went to mana")
		assert.Equal(t, match.HAND, inMana.Zone, "and a creature came back out")
	})

	t.Run("without the count nothing happens", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		passTurnToSelf(t, scn, player, opponent)
		summonWithOwnMana(t, scn, player, uncannyTurnipUID)
		require.NoError(t, scn.WaitForEventLoop())

		// Two Uncanny Turnips paid for the third, and without the wave striker
		// count nothing was added to or taken from the mana zone beyond that.
		mana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, mana, 2)
	})
}
