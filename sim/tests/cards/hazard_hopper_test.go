package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	hazardHopperUID = "999a432e-d38a-4f76-9bb7-d0dd54a3f9ad"
	hazardHopperSrc = "hazard_hopper_test_setup"
)

func TestHazardHopper(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		hopper := putCardInBattlezone(t, scn, player.Player, hazardHopperUID, hazardHopperSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, hopper, "Hazard Hopper", 5000, 4, []string{civ.Nature})
		assert.True(t, hopper.HasFamily(family.GiantInsect))
	})

	t.Run("it comes home after breaking a shield", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		hopper := putCardInBattlezone(t, scn, player.Player, hazardHopperUID, hazardHopperSrc)
		passTurnToSelf(t, scn, player, opponent)

		breakAShieldWith(t, scn, player, opponent, hopper)
		require.NoError(t, scn.ActionEndTurn(player))

		assert.Equal(t, match.HAND, hopper.Zone)
	})

	t.Run("two copies are tracked separately", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, player.Player, hazardHopperUID, hazardHopperSrc)
		bystander := putCardInBattlezone(t, scn, player.Player, hazardHopperUID, hazardHopperSrc)
		passTurnToSelf(t, scn, player, opponent)

		breakAShieldWith(t, scn, player, opponent, attacker)
		require.NoError(t, scn.ActionEndTurn(player))

		assert.Equal(t, match.HAND, attacker.Zone, "the one that attacked comes home")
		assert.Equal(t, match.BATTLEZONE, bystander.Zone, "the one that did not stays out")
	})
}
