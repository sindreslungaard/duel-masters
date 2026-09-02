package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	gigiosHammerUID      = "ca9c45d2-d280-4a57-9558-eb39a880cd3d"
	gigiosHammerAllyUID  = "a808b98c-2de7-412b-970c-a3b925bf43c2" // Deklowaz, the Terminator (Spirit Quartz, tap ability, no target needed)
	gigiosHammerOtherUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (Human, not Spirit Quartz)
	gigiosHammerSetupSrc = "gigios_hammer_test_setup"
)

// fx.ChooseAFamily offers family.GetAllFamilies() in order; Spirit Quartz is
// Deklowaz's race. Computed rather than hardcoded so a new family added
// anywhere in that list cannot silently shift this index out from under the test.
var gigiosHammerSpiritQuartzChoice = slices.Index(family.GetAllFamilies(), family.SpiritQuartz)

func TestGigiosHammer(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		hammer := putCardInBattlezone(t, scn, player.Player, gigiosHammerUID, gigiosHammerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, hammer, "Gigio's Hammer", 2000, 3, []string{civ.Fire})
		assert.True(t, hammer.HasFamily(family.Xenoparts))
		assert.True(t, hammer.HasCondition(cnd.TapAbility))
	})

	t.Run("the chosen race cannot use a tap ability instead of attacking", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		hammer := putCardInBattlezone(t, scn, player.Player, gigiosHammerUID, gigiosHammerSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, gigiosHammerAllyUID, gigiosHammerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)
		require.True(t, ally.HasCondition(cnd.TapAbility))

		require.NoError(t, scn.ActionUseTapAbility(player, hammer.ID))
		require.NoError(t, scn.SubmitChoice(player, gigiosHammerSpiritQuartzChoice))
		require.NoError(t, scn.WaitForEventLoop())
		assert.True(t, hammer.Tapped, "using the tap ability taps it")

		// The persistent effect it registers only takes effect from the next
		// event onward, which is why the grant and the tap-ability refusal are
		// both observed through this next action rather than beforehand.
		//
		// Official ruling (Slime Veil, the same "attacks if able" wording):
		// "On your turn, when you have the option either to attack with your
		// creature or to use its tap ability... you must attack with it."
		require.Error(t, scn.ActionUseTapAbility(player, ally.ID), "the tap ability should be refused")
		assert.False(t, ally.Tapped, "a refused ability does not tap the creature")
		require.True(t, ally.HasCondition(cnd.PowerAttacker))
		assert.Equal(t, 5000, scn.Match.GetPower(ally, false))
		assert.Equal(t, 9000, scn.Match.GetPower(ally, true), "5000 plus power attacker +4000")

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		_, err = scn.ActionAttackPlayer(player, ally.ID)
		require.NoError(t, err, "attacking with it instead is allowed")
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))

		require.NoError(t, scn.ActionEndTurn(player))
		assert.False(t, scn.Match.IsPlayerTurn(player.Player), "attacking satisfied the requirement")
		assert.False(t, ally.HasCondition(cnd.PowerAttacker), "the grant expired at the end of the turn")
	})

	t.Run("does not restrict creatures of a different race", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		hammer := putCardInBattlezone(t, scn, player.Player, gigiosHammerUID, gigiosHammerSetupSrc)
		other := putCardInBattlezone(t, scn, player.Player, gigiosHammerOtherUID, gigiosHammerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionUseTapAbility(player, hammer.ID))
		require.NoError(t, scn.SubmitChoice(player, gigiosHammerSpiritQuartzChoice))
		require.NoError(t, scn.WaitForEventLoop())

		require.NoError(t, scn.ActionEndTurn(player))
		assert.False(t, scn.Match.IsPlayerTurn(player.Player), "an untouched creature never holds up the turn")
		assert.False(t, other.HasCondition(cnd.PowerAttacker))
	})
}
