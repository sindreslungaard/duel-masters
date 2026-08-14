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
	lamielUID   = "e34eb3e3-f4b1-429a-90da-c5cf96a767da"
	lamielSetup = "lamiel_destiny_enforcer_test_setup"
)

func TestLamielDestinyEnforcer(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lamiel := putCardInBattlezone(t, scn, player.Player, lamielUID, lamielSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, lamiel, "Lamiel, Destiny Enforcer", 3000, 5, []string{civ.Light})
		assert.True(t, lamiel.HasFamily(family.Berserker))
		assert.True(t, lamiel.HasCondition(cnd.WaveStriker))
	})

	t.Run("it may draw when one of its controller's creatures dies on the opponent's turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, lamielUID, lamielSetup)
		addWaveStrikerFillers(t, scn, player, 2)
		victim := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, lamielSetup)

		attacker := putCardInBattlezone(t, scn, opponent.Player, waveStrikerBigCreatureUID, lamielSetup)

		passTurnToSelf(t, scn, player, opponent)
		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player))

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		// Killed in battle rather than by a direct Destroy call: the prompt has
		// to be raised by the event loop for the test goroutine to answer it.
		victim.Tapped = true
		require.NoError(t, scn.ActionAttackCreature(opponent, attacker.ID, victim.ID))

		// The prompt belongs to Lamiel's controller, who is not the player
		// taking the turn. It is raised while the destruction is still
		// resolving, so the creature has not reached the graveyard yet.
		answerInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, victim.Zone)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore)+1)
	})

	t.Run("it does not trigger on its controller's own turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, lamielUID, lamielSetup)
		addWaveStrikerFillers(t, scn, player, 2)
		victim := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, lamielSetup)

		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		scn.Match.Destroy(victim, victim, match.DestroyedByMiscAbility)
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
	})

	t.Run("without the count it offers nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, lamielUID, lamielSetup)
		victim := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, lamielSetup)

		passTurnToSelf(t, scn, player, opponent)
		require.NoError(t, scn.ActionEndTurn(player))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		scn.Match.Destroy(victim, victim, match.DestroyedByMiscAbility)
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
	})
}
