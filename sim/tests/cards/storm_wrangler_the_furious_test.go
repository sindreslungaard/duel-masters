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
	stormWranglerUID        = "1c877a61-6ec5-4c4c-87ee-8d4ab417b476"
	stormWranglerBlockerUID = "c7fec5e8-4e56-451b-a7b6-ad08680703a4" // La Byle, Seeker of the Winds (5000, Blocker)
	stormWranglerSetupSrc   = "storm_wrangler_test_setup"
)

func TestStormWranglerTheFurious(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, stormWranglerUID, stormWranglerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Storm Wrangler, the Furious", 5000, 4, []string{civ.Nature})
		assert.True(t, card.HasFamily(family.BeastFolk))
		assert.True(t, card.HasCondition(cnd.Evolution))
	})

	t.Run("gets +3000 power and wins when its chosen blocker blocks it", func(t *testing.T) {
		scn, player, _, attacker, blocker := setupStormWranglerAttack(t)

		attackStormWranglerForcingBlocker(t, scn, player, attacker, blocker)

		assert.Equal(t, match.BATTLEZONE, attacker.Zone, "8000 beats the blocker's 5000")
		assert.True(t, attacker.Tapped)
		assert.Equal(t, match.GRAVEYARD, blocker.Zone)
		assert.True(t, attacker.HasCondition(cnd.PowerAmplifier))
	})

	t.Run("the ability survives across turns, not just the turn it entered the battle zone", func(t *testing.T) {
		scn, player, opponent, attacker, _ := setupStormWranglerAttack(t)

		// A full round trip with no attack: an end-of-turn step for both
		// players passes before Storm Wrangler ever gets to attack.
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))
		require.True(t, scn.Match.IsPlayerTurn(player.Player))

		secondBlocker := putCardInBattlezone(t, scn, opponent.Player, stormWranglerBlockerUID, stormWranglerSetupSrc)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))
		require.True(t, scn.Match.IsPlayerTurn(player.Player))
		require.True(t, secondBlocker.HasCondition(cnd.Blocker))

		attackStormWranglerForcingBlocker(t, scn, player, attacker, secondBlocker)

		assert.Equal(t, match.BATTLEZONE, attacker.Zone, "the +3000 ability must still be active on a later turn")
		assert.Equal(t, match.GRAVEYARD, secondBlocker.Zone)
	})

	t.Run("the bonus expires at the end of the turn", func(t *testing.T) {
		scn, player, opponent, attacker, blocker := setupStormWranglerAttack(t)

		attackStormWranglerForcingBlocker(t, scn, player, attacker, blocker)
		require.True(t, attacker.HasCondition(cnd.PowerAmplifier))

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.False(t, attacker.HasCondition(cnd.PowerAmplifier), "the bonus is only \"until the end of the turn\"")
	})

	t.Run("an unblocked attack gets no bonus", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, player.Player, stormWranglerUID, stormWranglerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		action, err := scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID))

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount-1)
		assert.False(t, attacker.HasCondition(cnd.PowerAmplifier))
	})
}

// setupStormWranglerAttack puts Storm Wrangler in the battle zone for the
// current player and a 5000 power blocker in the opponent's, then cycles to a
// fresh turn so both have their intrinsic conditions.
func setupStormWranglerAttack(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference, *match.Card, *match.Card) {
	t.Helper()

	scn, player, opponent := setupDuel(t)
	attacker := putCardInBattlezone(t, scn, player.Player, stormWranglerUID, stormWranglerSetupSrc)
	blocker := putCardInBattlezone(t, scn, opponent.Player, stormWranglerBlockerUID, stormWranglerSetupSrc)
	passTurnToSelf(t, scn, player, opponent)

	require.True(t, blocker.HasCondition(cnd.Blocker))

	return scn, player, opponent, attacker, blocker
}

// attackStormWranglerForcingBlocker attacks the opposing player with attacker,
// breaks the first offered shield, then answers Storm Wrangler's own
// "choose a blocker to force" prompt with blocker and waits for the battle to
// resolve.
func attackStormWranglerForcingBlocker(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, attacker *match.Card, blocker *match.Card) {
	t.Helper()

	action, err := scn.ActionAttackPlayer(player, attacker.ID)
	require.NoError(t, err)

	messageStart, err := scn.MessageCount(player)
	require.NoError(t, err)
	require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID))

	_, err = scn.WaitForAction(player, messageStart)
	require.NoError(t, err, "expected Storm Wrangler's forced-blocker prompt to be open")
	require.NoError(t, scn.SubmitAction(player, blocker.ID))
	require.NoError(t, scn.WaitForEventLoop())
}
