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
	revivalSoldierUID = "ee2ae2d8-1ad3-4864-86ba-ce3c2fdb8ad5"
	revivalSetupSrc   = "revival_soldier_test_setup"
)

func TestRevivalSoldier(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		soldier := putCardInBattlezone(t, scn, player.Player, revivalSoldierUID, revivalSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, soldier, "Revival Soldier", 2000, 3, []string{civ.Water})
		assert.True(t, soldier.HasFamily(family.Merfolk))
		assert.True(t, soldier.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count it grows and comes back instead of dying", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		soldier := putCardInBattlezone(t, scn, player.Player, revivalSoldierUID, revivalSetupSrc)
		addWaveStrikerFillers(t, scn, player, 2)
		passTurnToSelf(t, scn, player, opponent)

		require.Equal(t, 6000, scn.Match.GetPower(soldier, false))

		scn.Match.Destroy(soldier, soldier, match.DestroyedByMiscAbility)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, soldier.Zone)
	})

	t.Run("without the count it simply dies", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		soldier := putCardInBattlezone(t, scn, player.Player, revivalSoldierUID, revivalSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		require.Equal(t, 2000, scn.Match.GetPower(soldier, false))

		scn.Match.Destroy(soldier, soldier, match.DestroyedByMiscAbility)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, soldier.Zone)
	})
}
