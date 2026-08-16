package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	slimeVeilUID       = "42949cc1-d2e4-459e-9ccf-ac93ebe636de"
	slimeVeilAllyUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000, no ability)
	slimeVeilTapperUID = "a808b98c-2de7-412b-970c-a3b925bf43c2" // Deklowaz, the Terminator (tap ability, no target needed)
	slimeVeilSetupSrc  = "slime_veil_test_setup"
)

func TestSlimeVeil(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(slimeVeilUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Slime Veil", spell.Name)
		assert.Equal(t, 1, spell.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, spell.Civs)
		assert.Equal(t, []string{civ.Darkness}, spell.ManaRequirement)
	})

	t.Run("forces the opponent's creatures to attack before ending their turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		ally := putCardInBattlezone(t, scn, opponent.Player, slimeVeilAllyUID, slimeVeilSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, slimeVeilUID)

		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player), "it is now the opponent's forced turn")

		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, scn.Match.IsPlayerTurn(opponent.Player), "the turn cannot end while Slime Veil's target can still attack")

		// Tapping stands in for the ally having attacked.
		ally.Tapped = true
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, scn.Match.IsPlayerTurn(player.Player))
	})

	t.Run("the opponent cannot use a tap ability instead of attacking", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		tapper := putCardInBattlezone(t, scn, opponent.Player, slimeVeilTapperUID, slimeVeilSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, slimeVeilUID)
		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player))

		// Official ruling: "On your turn, when you have the option either to
		// attack with your creature or to use its tap ability, [Slime Veil]
		// says that you must attack with it."
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
		tapper := putCardInBattlezone(t, scn, opponent.Player, slimeVeilTapperUID, slimeVeilSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, slimeVeilUID)
		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player))

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		_, err = scn.ActionAttackPlayer(opponent, tapper.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))
		require.NoError(t, scn.ActionEndTurn(opponent))
		require.True(t, scn.Match.IsPlayerTurn(player.Player), "Slime Veil's effect has now expired")

		require.NoError(t, scn.ActionEndTurn(player))
		require.False(t, scn.Match.IsPlayerTurn(player.Player), "the opponent's second turn is unrestricted")

		require.NoError(t, scn.ActionUseTapAbility(opponent, tapper.ID))
		require.NoError(t, scn.WaitForEventLoop())
		assert.True(t, tapper.Tapped, "the tap ability resolved")
	})

	t.Run("does not restrict the caster's own creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		tapper := putCardInBattlezone(t, scn, player.Player, slimeVeilTapperUID, slimeVeilSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, slimeVeilUID)

		require.NoError(t, scn.ActionUseTapAbility(player, tapper.ID), "Slime Veil only forces the opponent's creatures")
		require.NoError(t, scn.WaitForEventLoop())
		assert.True(t, tapper.Tapped)
	})
}
