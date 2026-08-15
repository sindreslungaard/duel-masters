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
	tauntingSkyterrorUID       = "93ac0e42-2ab3-4e1c-b6f4-ac8caaa87f88"
	tauntingSkyterrorAllyUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	tauntingSkyterrorTapperUID = "a808b98c-2de7-412b-970c-a3b925bf43c2" // Deklowaz, the Terminator (tap ability, no target needed)
	tauntingSkyterrorSetupSrc  = "taunting_skyterror_test_setup"
)

func TestTauntingSkyterror(t *testing.T) {
	t.Run("holds up the opponent's turn while it is tapped", func(t *testing.T) {
		scn, owner, opponent := setupTauntingSkyterrorTest(t)
		skyterror := putCardInBattlezone(t, scn, owner.Player, tauntingSkyterrorUID, tauntingSkyterrorSetupSrc)
		ally := putCardInBattlezone(t, scn, opponent.Player, tauntingSkyterrorAllyUID, tauntingSkyterrorSetupSrc)

		assert.Equal(t, "Taunting Skyterror", skyterror.Name)
		assert.Equal(t, 3000, skyterror.Power)
		assert.Equal(t, 5, skyterror.ManaCost)
		assert.Equal(t, []string{civ.Fire}, skyterror.Civs)
		assert.True(t, skyterror.HasFamily(family.ArmoredWyvern))

		require.NoError(t, scn.ActionEndTurn(owner))
		skyterror.Tapped = true

		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.False(t, scn.Match.IsPlayerTurn(owner.Player), "the opponent still has an untapped creature")

		ally.Tapped = true
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, scn.Match.IsPlayerTurn(owner.Player))
	})

	t.Run("does nothing while it is untapped", func(t *testing.T) {
		scn, owner, opponent := setupTauntingSkyterrorTest(t)
		putCardInBattlezone(t, scn, owner.Player, tauntingSkyterrorUID, tauntingSkyterrorSetupSrc)
		putCardInBattlezone(t, scn, opponent.Player, tauntingSkyterrorAllyUID, tauntingSkyterrorSetupSrc)

		require.NoError(t, scn.ActionEndTurn(owner))
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, scn.Match.IsPlayerTurn(owner.Player))
	})

	t.Run("the opponent cannot use a tap ability instead of attacking", func(t *testing.T) {
		scn, owner, opponent := setupTauntingSkyterrorTest(t)
		skyterror := putCardInBattlezone(t, scn, owner.Player, tauntingSkyterrorUID, tauntingSkyterrorSetupSrc)
		tapper := putCardInBattlezone(t, scn, opponent.Player, tauntingSkyterrorTapperUID, tauntingSkyterrorSetupSrc)

		require.NoError(t, scn.ActionEndTurn(owner))
		skyterror.Tapped = true
		require.False(t, scn.Match.IsPlayerTurn(owner.Player))
		require.True(t, tapper.HasCondition(cnd.TapAbility), "the opponent's own untap step just ran")

		// Official ruling (Slime Veil, the same "attacks if able" wording):
		// "when you have the option either to attack with your creature or to
		// use its tap ability... you must attack with it."
		require.Error(t, scn.ActionUseTapAbility(opponent, tapper.ID), "the tap ability should be refused")
		assert.False(t, tapper.Tapped, "a refused ability does not tap the creature")

		shields, err := owner.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		_, err = scn.ActionAttackPlayer(opponent, tapper.ID)
		require.NoError(t, err, "attacking with it instead is allowed")
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))

		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, scn.Match.IsPlayerTurn(owner.Player), "attacking satisfied the requirement")
	})
}

func setupTauntingSkyterrorTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	owner := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))

	return scn, owner, opponent
}
