package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	wingeyeMothUID    = "02f11a2f-ff50-4e0d-80bd-2995be9d3dd8"
	wingeyeSmallUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	wingeyeBiggerUID  = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	wingeyeBiggestUID = "43abeec5-0597-43b3-93cf-766b95d19b5b" // Forest Hornet (4000)
	wingeyeSrc        = "wingeye_moth_test_setup"
)

func TestWingeyeMoth(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, wingeyeMothUID, wingeyeSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Wingeye Moth", 3000, 5, []string{civ.Nature})
		assert.True(t, card.HasFamily(family.GiantInsect))
	})

	t.Run("an empty board opposite is nothing to beat", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, wingeyeMothUID, wingeyeSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(hand)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(opponent))

		_, err = scn.WaitForAction(player, promptStart)
		require.NoError(t, err, "the extra draw should have been offered")
		require.NoError(t, scn.SubmitAction(player))
		settleTurn(t, scn)

		hand, err = player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount+2, "the turn's own draw plus the extra one")
	})

	t.Run("outmuscling everything opposite offers the draw", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, wingeyeMothUID, wingeyeSrc)
		putCardInBattlezone(t, scn, opponent.Player, wingeyeSmallUID, wingeyeSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(hand)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(opponent))

		_, err = scn.WaitForAction(player, promptStart)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player))
		settleTurn(t, scn)

		hand, err = player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount+2)
	})

	t.Run("the extra draw may be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, wingeyeMothUID, wingeyeSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(hand)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(opponent))

		_, err = scn.WaitForAction(player, promptStart)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		settleTurn(t, scn)

		hand, err = player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount+1, "only the turn's own draw")
	})

	t.Run("a creature it cannot beat keeps it quiet", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, wingeyeMothUID, wingeyeSrc)
		putCardInBattlezone(t, scn, opponent.Player, wingeyeBiggerUID, wingeyeSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(hand)

		// No prompt, so the turn runs all the way through on its own.
		require.NoError(t, scn.ActionEndTurn(opponent))
		settleTurn(t, scn)

		hand, err = player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount+1)
	})

	t.Run("any creature of yours can be the one that qualifies", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, wingeyeMothUID, wingeyeSrc)
		putCardInBattlezone(t, scn, player.Player, wingeyeBiggestUID, wingeyeSrc)
		putCardInBattlezone(t, scn, opponent.Player, wingeyeSmallUID, wingeyeSrc)

		// The moth itself is out-powered by nothing here, but the check is over
		// every creature its controller has, not only over itself.
		require.NoError(t, scn.ActionEndTurn(player))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(opponent))

		_, err = scn.WaitForAction(player, promptStart)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		settleTurn(t, scn)
	})

	t.Run("it stays quiet on the opponent's draw step", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, wingeyeMothUID, wingeyeSrc)

		// Handing the turn over runs the opponent's draw step, which is not
		// this creature's.
		require.NoError(t, scn.ActionEndTurn(player))
		settleTurn(t, scn)

		assert.True(t, scn.Match.IsPlayerTurn(opponent.Player))
	})
}
