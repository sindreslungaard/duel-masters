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
	// Apocalypse Vise (fire, cost 7): destroy any number of the opponent's
	// creatures with total power 8000 or less, chosen freely and rechecked
	// against the total until it fits.
	apocalypseViseUID = "e0558aef-d2d3-4111-aa79-965cdc604f57"
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

	t.Run("Apocalypse Vise recomputes it using the post-cast hand size", func(t *testing.T) {
		// Reported bug: "if opponent has a scarab and I have a hand of four
		// cards and cast Apoc, the pincer should be able to be targeted and
		// destroyed, but Shobu treats it as still having 9k power (acts as
		// though the Apoc is still in hand)."
		scn, player, opponent := setupDuel(t)
		scarab := putCardInBattlezone(t, scn, opponent.Player, pincerScarabUID, pincerScarabSrc)

		// A hand of exactly 4 cards, one of which is Apocalypse Vise itself:
		// once it is cast, 3 remain, putting Scarab at 1000+2000*3=7000 -
		// destructible alone, where the stale 4-card count would keep it at
		// 9000 and out of reach.
		emptyHand(t, player, pincerScarabSrc)
		for range 3 {
			_, err := player.Player.SpawnCard(pincerScarabFillUID, match.HAND)
			require.NoError(t, err)
		}

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, apocalypseViseUID)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err, "expected the destroy-selection prompt to be open")
		offeredCardIDs := make([]string, 0, len(action.Cards))
		for _, offered := range action.Cards {
			offeredCardIDs = append(offeredCardIDs, offered.CardID)
		}
		require.Contains(t, offeredCardIDs, scarab.ID)

		require.NoError(t, scn.SubmitAction(player, scarab.ID))
		require.NoError(t, scn.WaitForEventLoop(), "Scarab alone should have been an acceptable selection, not rejected for exceeding 8000 total power")

		assert.Equal(t, match.GRAVEYARD, scarab.Zone)
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})
}
