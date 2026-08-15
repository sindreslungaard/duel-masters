package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	forcedFrenzyUID       = "fb0fe721-bc21-4f10-a807-7aa698c8e38a"
	forcedFrenzyAllyUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000, no ability)
	forcedFrenzyTapperUID = "a808b98c-2de7-412b-970c-a3b925bf43c2" // Deklowaz, the Terminator (tap ability, no target needed)
	forcedFrenzySetupSrc  = "forced_frenzy_test_setup"
)

func TestForcedFrenzy(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		spell, err := player.Player.SpawnCard(forcedFrenzyUID, match.HAND)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Forced Frenzy", spell.Name)
		assert.Equal(t, 3, spell.ManaCost)
		assert.Equal(t, []string{civ.Fire}, spell.Civs)
		assert.Equal(t, []string{civ.Fire}, spell.ManaRequirement)
		assert.True(t, spell.HasCondition(cnd.ShieldTrigger))
	})

	t.Run("forces the opponent's creatures to attack before ending their turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		ally := putCardInBattlezone(t, scn, opponent.Player, forcedFrenzyAllyUID, forcedFrenzySetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, forcedFrenzyUID)

		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player), "it is now the opponent's forced turn")

		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, scn.Match.IsPlayerTurn(opponent.Player), "the turn cannot end while Forced Frenzy's target can still attack")

		// Tapping stands in for the ally having attacked.
		ally.Tapped = true
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, scn.Match.IsPlayerTurn(player.Player))
	})

	t.Run("the opponent cannot use a tap ability instead of attacking", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		tapper := putCardInBattlezone(t, scn, opponent.Player, forcedFrenzyTapperUID, forcedFrenzySetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, forcedFrenzyUID)
		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player))

		// Official ruling (Slime Veil, the same "attacks if able" wording):
		// "On your turn, when you have the option either to attack with your
		// creature or to use its tap ability... you must attack with it."
		require.Error(t, scn.ActionUseTapAbility(opponent, tapper.ID), "the tap ability should be refused")
		assert.False(t, tapper.Tapped, "a refused ability does not tap the creature")

		// Attacking with it instead is allowed, and satisfies the requirement.
		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		_, err = scn.ActionAttackPlayer(opponent, tapper.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))

		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, scn.Match.IsPlayerTurn(player.Player), "attacking satisfied the requirement")
	})

	t.Run("the tap ability works normally again once the forced turn is over", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		tapper := putCardInBattlezone(t, scn, opponent.Player, forcedFrenzyTapperUID, forcedFrenzySetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, forcedFrenzyUID)
		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player))

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		_, err = scn.ActionAttackPlayer(opponent, tapper.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))
		require.NoError(t, scn.ActionEndTurn(opponent))
		require.True(t, scn.Match.IsPlayerTurn(player.Player), "Forced Frenzy's effect has now expired")

		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player), "the opponent's second turn is unrestricted")

		require.NoError(t, scn.ActionUseTapAbility(opponent, tapper.ID))
		require.NoError(t, scn.WaitForEventLoop())
		assert.True(t, tapper.Tapped, "the tap ability resolved")
	})

	t.Run("does not restrict the caster's own creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		tapper := putCardInBattlezone(t, scn, player.Player, forcedFrenzyTapperUID, forcedFrenzySetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, forcedFrenzyUID)

		require.NoError(t, scn.ActionUseTapAbility(player, tapper.ID), "Forced Frenzy only forces the opponent's creatures")
		require.NoError(t, scn.WaitForEventLoop())
		assert.True(t, tapper.Tapped)
	})
}
