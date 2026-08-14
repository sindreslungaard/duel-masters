package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	luckyBallUID      = "71b80172-8b18-4140-99bf-d8bc4201e07e"
	luckyBallSetupSrc = "lucky_ball_test_setup"
)

func TestLuckyBall(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		ball := putCardInBattlezone(t, scn, player.Player, luckyBallUID, luckyBallSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, ball, "Lucky Ball", 3000, 4, []string{civ.Water})
		assert.True(t, ball.HasFamily(family.CyberVirus))
	})

	t.Run("it draws when the opponent is down to three shields", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		trimOpponentShields(t, opponent, 3, luckyBallSetupSrc)

		handBefore := summonLuckyBall(t, scn, player)

		answerDrawUpTo(t, scn, player, 2, true)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		// Lucky Ball itself left the hand to be summoned, and two cards came
		// back in its place.
		assert.Len(t, hand, handBefore+1)
	})

	t.Run("it does nothing while the opponent still has four shields", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		trimOpponentShields(t, opponent, 4, luckyBallSetupSrc)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		summonLuckyBall(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment should ask anything")
	})
}

func summonLuckyBall(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) int {
	t.Helper()

	ball, err := player.Player.SpawnCard(luckyBallUID, match.HAND)
	require.NoError(t, err)
	for range 4 {
		_, err := player.Player.SpawnCard(luckyBallUID, match.MANAZONE)
		require.NoError(t, err)
	}

	hand, err := player.Player.Container(match.HAND)
	require.NoError(t, err)
	before := len(hand)

	require.NoError(t, scn.ActionPlayCard(player, ball.ID))

	return before
}

func trimOpponentShields(t *testing.T, opponent *match.PlayerReference, keep int, source string) {
	t.Helper()

	shields, err := opponent.Player.Container(match.SHIELDZONE)
	require.NoError(t, err)

	for _, shield := range shields[keep:] {
		_, err := opponent.Player.MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, source)
		require.NoError(t, err)
	}
}
