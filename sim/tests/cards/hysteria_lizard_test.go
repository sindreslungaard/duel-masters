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
	hysteriaLizardUID = "e9edcaed-0434-4e97-b058-bc4fe955ac08"
	hysteriaSetupSrc  = "hysteria_lizard_test_setup"
)

func TestHysteriaLizard(t *testing.T) {
	t.Run("printed characteristics and power attacker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// Held in hand for this one: on the board it refuses to let the turn
		// end, which is what the turn handover below would try to do.
		lizard, err := player.Player.SpawnCard(hysteriaLizardUID, match.HAND)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, lizard, "Hysteria Lizard", 3000, 4, []string{civ.Fire})
		assert.True(t, lizard.HasFamily(family.MeltWarrior))
		assert.True(t, lizard.HasCondition(cnd.PowerAttacker))

		assert.Equal(t, 3000, scn.Match.GetPower(lizard, false))
		assert.Equal(t, 6000, scn.Match.GetPower(lizard, true), "power attacker +3000 while attacking")
	})

	t.Run("the turn cannot be ended while it is able to attack", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		lizard := putCardInBattlezone(t, scn, player.Player, hysteriaLizardUID, hysteriaSetupSrc)

		warningStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(player))

		warnings, err := scn.Warnings(player, warningStart)
		require.NoError(t, err)
		assert.NotEmpty(t, warnings, "ending the turn should be refused")
		assert.True(t, scn.Match.IsPlayerTurn(player.Player), "the turn did not end")
		assert.Equal(t, false, lizard.Tapped)
	})

	t.Run("a tapped creature no longer holds up the turn", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		lizard := putCardInBattlezone(t, scn, player.Player, hysteriaLizardUID, hysteriaSetupSrc)

		lizard.Tapped = true
		require.NoError(t, scn.ActionEndTurn(player))

		assert.False(t, scn.Match.IsPlayerTurn(player.Player))
	})
}
