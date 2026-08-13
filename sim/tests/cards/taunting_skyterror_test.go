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
	tauntingSkyterrorUID      = "93ac0e42-2ab3-4e1c-b6f4-ac8caaa87f88"
	tauntingSkyterrorAllyUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	tauntingSkyterrorSetupSrc = "taunting_skyterror_test_setup"
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
}

func setupTauntingSkyterrorTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	owner := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))

	return scn, owner, opponent
}
