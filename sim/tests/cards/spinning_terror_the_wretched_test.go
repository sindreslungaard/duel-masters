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
	spinningTerrorUID      = "5bcab12f-5a17-4ade-b938-1c08a6290047"
	spinningTerrorSetupSrc = "spinning_terror_the_wretched_test_setup"
)

func TestSpinningTerrorTheWretched(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		terror := putCardInBattlezone(t, scn, player.Player, spinningTerrorUID, spinningTerrorSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, terror, "Spinning Terror, the Wretched", 1000, 2, []string{civ.Darkness})
		assert.True(t, terror.HasFamily(family.DevilMask))
		assert.Equal(t, 1000, scn.Match.GetPower(terror, false), "no tapped opposing creatures yet")
	})

	t.Run("it grows with each tapped creature the opponent has", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		terror := putCardInBattlezone(t, scn, player.Player, spinningTerrorUID, spinningTerrorSetupSrc)

		first := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, spinningTerrorSetupSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, spinningTerrorSetupSrc)
		own := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, spinningTerrorSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 1000, scn.Match.GetPower(terror, false))

		first.Tapped = true
		assert.Equal(t, 3000, scn.Match.GetPower(terror, false))

		second.Tapped = true
		assert.Equal(t, 5000, scn.Match.GetPower(terror, false))

		own.Tapped = true
		assert.Equal(t, 5000, scn.Match.GetPower(terror, false), "its own tapped creatures do not count")

		first.Tapped = false
		assert.Equal(t, 3000, scn.Match.GetPower(terror, false), "the bonus follows the board")
	})

	t.Run("it counts nothing from outside the battle zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		terror := putCardInBattlezone(t, scn, player.Player, spinningTerrorUID, spinningTerrorSetupSrc)
		tapped := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, spinningTerrorSetupSrc)
		tapped.Tapped = true

		passTurnToSelf(t, scn, player, opponent)
		tapped.Tapped = true
		require.Equal(t, 3000, scn.Match.GetPower(terror, false))

		_, err := player.Player.MoveCard(terror.ID, match.BATTLEZONE, match.HAND, spinningTerrorSetupSrc)
		require.NoError(t, err)

		assert.Equal(t, 1000, scn.Match.GetPower(terror, false))
	})
}
