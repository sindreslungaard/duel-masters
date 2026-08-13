package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	pincerScarabUID     = "7df9100a-a893-45ea-ab1a-8312b4232b65"
	pincerScarabFillUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	pincerScarabSrc     = "pincer_scarab_test_setup"
)

func TestPincerScarab(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, pincerScarabUID, pincerScarabSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Pincer Scarab", 1000, 4, []string{civ.Nature})
		assert.True(t, card.HasFamily(family.GiantInsect))
	})

	t.Run("power tracks the opponent's hand", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		scarab := putCardInBattlezone(t, scn, player.Player, pincerScarabUID, pincerScarabSrc)

		emptyHand(t, opponent, pincerScarabSrc)
		assert.Equal(t, 1000, scn.Match.GetPower(scarab, false), "an empty hand is worth nothing")

		for range 2 {
			_, err := opponent.Player.SpawnCard(pincerScarabFillUID, match.HAND)
			require.NoError(t, err)
		}
		assert.Equal(t, 5000, scn.Match.GetPower(scarab, false))

		_, err := opponent.Player.SpawnCard(pincerScarabFillUID, match.HAND)
		require.NoError(t, err)
		assert.Equal(t, 7000, scn.Match.GetPower(scarab, false))

		// Its own hand is not what feeds it.
		_, err = player.Player.SpawnCard(pincerScarabFillUID, match.HAND)
		require.NoError(t, err)
		assert.Equal(t, 7000, scn.Match.GetPower(scarab, false))
	})

	t.Run("a full opposing hand gives it double breaker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		scarab := putCardInBattlezone(t, scn, player.Player, pincerScarabUID, pincerScarabSrc)

		emptyHand(t, opponent, pincerScarabSrc)
		for range 3 {
			_, err := opponent.Player.SpawnCard(pincerScarabFillUID, match.HAND)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		require.GreaterOrEqual(t, scn.Match.GetPower(scarab, false), 6000)
		assert.True(t, scarab.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("a thin opposing hand leaves it a single breaker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		scarab := putCardInBattlezone(t, scn, player.Player, pincerScarabUID, pincerScarabSrc)

		emptyHand(t, opponent, pincerScarabSrc)

		// The opponent draws once on the turn that passes through, so they end
		// up holding a single card.
		passTurnToSelf(t, scn, player, opponent)

		require.Less(t, scn.Match.GetPower(scarab, false), 6000)
		assert.False(t, scarab.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("the breaker follows any bonus, not only the hand", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		scarab := putCardInBattlezone(t, scn, player.Player, pincerScarabUID, pincerScarabSrc)

		emptyHand(t, opponent, pincerScarabSrc)
		passTurnToSelf(t, scn, player, opponent)
		require.False(t, scarab.HasCondition(cnd.DoubleBreaker))

		// Granted after the turn transition, which is where conditions are
		// cleared, and small enough that only the hand on top of it clears 6000.
		scarab.AddUniqueSourceCondition(cnd.PowerAmplifier, 5000, pincerScarabSrc)
		require.Equal(t, 8000, scn.Match.GetPower(scarab, false))

		// The tier is re-checked on the next event of any kind.
		emptyHand(t, player, pincerScarabSrc)

		assert.True(t, scarab.HasCondition(cnd.DoubleBreaker))
	})
}
