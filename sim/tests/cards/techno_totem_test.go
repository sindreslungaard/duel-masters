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
	technoTotemUID      = "a1ef4e4e-d8c9-4b33-bfb0-3f8f7d021ff5"
	technoTotemAllyUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	technoTotemSetupSrc = "techno_totem_test_setup"
)

func TestTechnoTotem(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupTechnoTotemTest(t)
		totem := putCardInBattlezone(t, scn, player.Player, technoTotemUID, technoTotemSetupSrc)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.Equal(t, "Techno Totem", totem.Name)
		assert.Equal(t, 5000, totem.Power)
		assert.Equal(t, 4, totem.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Nature}, totem.Civs)
		assert.True(t, totem.HasFamily(family.MysteryTotem))
		assert.True(t, totem.HasCondition(cnd.TapAbility))
	})

	t.Run("grants power attacker to other own creatures only while tapped", func(t *testing.T) {
		scn, player, opponent := setupTechnoTotemTest(t)
		totem := putCardInBattlezone(t, scn, player.Player, technoTotemUID, technoTotemSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, technoTotemAllyUID, technoTotemSetupSrc)
		enemy := putCardInBattlezone(t, scn, opponent.Player, technoTotemAllyUID, technoTotemSetupSrc)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.Equal(t, 2000, scn.Match.GetPower(ally, true), "untapped, it grants nothing")

		totem.Tapped = true
		require.NoError(t, scn.ActionEndTurn(player))

		assert.Equal(t, 3500, scn.Match.GetPower(ally, true), "power attacker +1500 while attacking")
		assert.Equal(t, 2000, scn.Match.GetPower(ally, false), "and nothing while not attacking")
		assert.Equal(t, 2000, scn.Match.GetPower(enemy, true), "the opponent's creatures are unaffected")
		assert.Equal(t, 5000, scn.Match.GetPower(totem, true), "it does not boost itself")

		totem.Tapped = false
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.Equal(t, 2000, scn.Match.GetPower(ally, true), "the grant follows the tapped state")
	})

	t.Run("stops granting when it leaves the battle zone", func(t *testing.T) {
		scn, player, opponent := setupTechnoTotemTest(t)
		totem := putCardInBattlezone(t, scn, player.Player, technoTotemUID, technoTotemSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, technoTotemAllyUID, technoTotemSetupSrc)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		totem.Tapped = true
		require.NoError(t, scn.ActionEndTurn(player))
		require.Equal(t, 3500, scn.Match.GetPower(ally, true))

		moved, err := player.Player.MoveCard(totem.ID, match.BATTLEZONE, match.GRAVEYARD, technoTotemSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		assert.Equal(t, 2000, scn.Match.GetPower(ally, true))
		assert.False(t, ally.HasCondition(cnd.PowerAttacker))
	})

	t.Run("taps an opposing creature with its tap ability", func(t *testing.T) {
		scn, player, opponent := setupTechnoTotemTest(t)
		totem := putCardInBattlezone(t, scn, player.Player, technoTotemUID, technoTotemSetupSrc)
		enemy := putCardInBattlezone(t, scn, opponent.Player, technoTotemAllyUID, technoTotemSetupSrc)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		require.False(t, enemy.Tapped)
		require.NoError(t, scn.ActionUseTapAbility(player, totem.ID))

		assert.True(t, enemy.Tapped, "the only legal target is tapped without a prompt")
		assert.True(t, totem.Tapped, "using the tap ability taps this creature")
	})
}

func setupTechnoTotemTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}
