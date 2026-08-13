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
	evilIncarnateUID  = "b888a7a8-d8c0-458a-8cf6-0b675ab9123e"
	evilIncarnateBase = "5bcab12f-5a17-4ade-b938-1c08a6290047" // Spinning Terror, the Wretched (Devil Mask)
	evilFodderUID     = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	evilIncarnateSrc  = "evil_incarnate_test_setup"
)

func TestEvilIncarnate(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		evil := putCardInBattlezone(t, scn, player.Player, evilIncarnateUID, evilIncarnateSrc)
		fodder := putCardInBattlezone(t, scn, player.Player, evilFodderUID, evilIncarnateSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		// The opponent gives one up at the start of their turn.
		answerEvilIncarnate(t, scn, opponent)

		assertPrinted(t, evil, "Evil Incarnate", 11000, 6, []string{civ.Darkness})
		assert.True(t, evil.HasFamily(family.DevilMask))
		assert.True(t, evil.HasCondition(cnd.Evolution))
		assert.True(t, evil.HasCondition(cnd.DoubleBreaker))
		assert.Equal(t, match.BATTLEZONE, fodder.Zone, "its controller's turn has not come round yet")
	})

	t.Run("the player whose turn it is gives up a creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, evilIncarnateUID, evilIncarnateSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, evilFodderUID, evilIncarnateSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		answerEvilIncarnate(t, scn, opponent)

		// The opponent had exactly one creature, so the choice resolved itself.
		assert.Equal(t, match.GRAVEYARD, theirs.Zone)
	})

	t.Run("its own controller pays too", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, evilIncarnateUID, evilIncarnateSrc)
		mine := putCardInBattlezone(t, scn, player.Player, evilFodderUID, evilIncarnateSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		answerEvilIncarnate(t, scn, opponent)

		require.NoError(t, scn.ActionEndTurn(opponent))
		answerEvilIncarnate(t, scn, player, mine.ID)

		assert.Equal(t, match.GRAVEYARD, mine.Zone, "\"each player's turn\" includes its own")
	})

	t.Run("a player with nothing on the board loses nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, evilIncarnateUID, evilIncarnateSrc)

		promptStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(opponent, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action", "an empty battle zone offers no choice")
	})

	t.Run("it stops once it has left the battle zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		evil := putCardInBattlezone(t, scn, player.Player, evilIncarnateUID, evilIncarnateSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, evilFodderUID, evilIncarnateSrc)

		_, err := player.Player.MoveCard(evil.ID, match.BATTLEZONE, match.GRAVEYARD, evilIncarnateSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		promptStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(opponent, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
		assert.Equal(t, match.BATTLEZONE, theirs.Zone)
	})

	t.Run("it evolves onto a Devil Mask", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		base := putCardInBattlezone(t, scn, player.Player, evilIncarnateBase, evilIncarnateSrc)

		evil := spawnForLater(t, player, evilIncarnateUID)
		for range 6 {
			_, err := player.Player.SpawnCard(evilIncarnateUID, match.MANAZONE)
			require.NoError(t, err)
		}

		// Summoned on its controller's own turn, before its start of turn
		// trigger has had a chance to fire at all.
		require.NoError(t, scn.ActionPlayCard(player, evil.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, evil.Zone)
		assert.Equal(t, match.HIDDENZONE, base.Zone, "the base goes under the evolution")
	})
}

// answerEvilIncarnate answers the start of turn prompt, if one was opened.
func answerEvilIncarnate(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, cardIDs ...string) {
	t.Helper()

	headers, err := scn.MessageHeaders(player, 0)
	require.NoError(t, err)
	if countHeaders(headers, "action") == 0 {
		return
	}

	answerInTurn(t, scn, player, cardIDs...)
	require.NoError(t, scn.WaitForEventLoop())
}
