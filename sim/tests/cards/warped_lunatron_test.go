package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	warpedLunatronUID = "fccd4a14-6b55-42bf-8dac-8ac311fc4571"
	warpedCreatureUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	warpedSetupSrc    = "warped_lunatron_test_setup"
)

func TestWarpedLunatron(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, warpedLunatronUID, warpedSetupSrc)

		passTurnThroughWarp(t, scn, player, opponent)

		assertPrinted(t, lunatron, "Warped Lunatron", 6000, 7, []string{civ.Water})
		assert.True(t, lunatron.HasFamily(family.CyberMoon))
		assert.True(t, lunatron.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("creatures stay tapped through the untap step", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, warpedLunatronUID, warpedSetupSrc)
		mine := putCardInBattlezone(t, scn, player.Player, warpedCreatureUID, warpedSetupSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, warpedCreatureUID, warpedSetupSrc)

		mine.Tapped = true
		theirs.Tapped = true

		passTurnThroughWarp(t, scn, player, opponent)

		assert.True(t, mine.Tapped, "the untap step no longer frees it")
		assert.True(t, theirs.Tapped, "the printed text covers both players' creatures")
	})

	t.Run("two cards of mana buy one untap", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, warpedLunatronUID, warpedSetupSrc)
		mine := putCardInBattlezone(t, scn, player.Player, warpedCreatureUID, warpedSetupSrc)
		other := putCardInBattlezone(t, scn, player.Player, warpedCreatureUID, warpedSetupSrc)
		mine.Tapped = true
		other.Tapped = true

		mana := make([]*match.Card, 0, 2)
		for range 2 {
			card, err := player.Player.SpawnCard(warpedCreatureUID, match.MANAZONE)
			require.NoError(t, err)
			mana = append(mana, card)
		}

		require.NoError(t, scn.ActionEndTurn(player))
		answerWarpIfOffered(t, scn, opponent)
		require.NoError(t, scn.ActionEndTurn(opponent))

		// The caster's own untap step: pay two mana, then say which single
		// creature that buys back.
		answerInTurn(t, scn, player, mana[0].ID, mana[1].ID)
		answerInTurn(t, scn, player, mine.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.False(t, mine.Tapped, "the trade untapped the one chosen")
		assert.True(t, other.Tapped, "two mana only buys one untap")
		assert.True(t, mana[0].Tapped, "and the mana was spent")
		assert.True(t, mana[1].Tapped)
	})

	t.Run("the trade can be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, warpedLunatronUID, warpedSetupSrc)
		mine := putCardInBattlezone(t, scn, player.Player, warpedCreatureUID, warpedSetupSrc)
		mine.Tapped = true

		for range 2 {
			_, err := player.Player.SpawnCard(warpedCreatureUID, match.MANAZONE)
			require.NoError(t, err)
		}

		require.NoError(t, scn.ActionEndTurn(player))
		answerWarpIfOffered(t, scn, opponent)
		require.NoError(t, scn.ActionEndTurn(opponent))

		cancelInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, mine.Tapped, "declining leaves it tapped")
	})

	t.Run("nothing is offered without the mana to pay", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, warpedLunatronUID, warpedSetupSrc)
		mine := putCardInBattlezone(t, scn, player.Player, warpedCreatureUID, warpedSetupSrc)
		mine.Tapped = true

		_, err := player.Player.SpawnCard(warpedCreatureUID, match.MANAZONE)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(player))
		answerWarpIfOffered(t, scn, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(opponent))
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action", "one card of mana is not enough to trade")
		assert.True(t, mine.Tapped)
	})

	t.Run("creatures untap normally once it has left", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, warpedLunatronUID, warpedSetupSrc)
		mine := putCardInBattlezone(t, scn, player.Player, warpedCreatureUID, warpedSetupSrc)
		mine.Tapped = true

		_, err := player.Player.MoveCard(lunatron.ID, match.BATTLEZONE, match.GRAVEYARD, warpedSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		passTurnToSelf(t, scn, player, opponent)

		assert.False(t, mine.Tapped)
		assert.False(t, mine.HasCondition(cnd.DoesntUntap), "the condition goes with its source")
	})
}

// passTurnThroughWarp hands the turn over and back, answering the trade the
// Lunatron offers each player at their untap step.
func passTurnThroughWarp(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, opponent *match.PlayerReference) {
	t.Helper()

	require.NoError(t, scn.ActionEndTurn(player))
	answerWarpIfOffered(t, scn, opponent)
	require.NoError(t, scn.ActionEndTurn(opponent))
	answerWarpIfOffered(t, scn, player)
}

// answerWarpIfOffered declines the trade when it was offered at all.
func answerWarpIfOffered(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) {
	t.Helper()

	headers, err := scn.MessageHeaders(player, 0)
	require.NoError(t, err)
	if countHeaders(headers, "action") == 0 {
		return
	}

	cancelInTurn(t, scn, player)
	require.NoError(t, scn.WaitForEventLoop())
}
