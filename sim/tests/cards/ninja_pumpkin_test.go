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
	ninjaPumpkinUID = "0b99b1d5-240b-4a4a-b505-e446df330c40"
	ninjaSetupSrc   = "ninja_pumpkin_test_setup"
)

func TestNinjaPumpkin(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		pumpkin := putCardInBattlezone(t, scn, player.Player, ninjaPumpkinUID, ninjaSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, pumpkin, "Ninja Pumpkin", 2000, 3, []string{civ.Nature})
		assert.True(t, pumpkin.HasFamily(family.WildVeggies))
		assert.True(t, pumpkin.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count it grows", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		pumpkin := putCardInBattlezone(t, scn, player.Player, ninjaPumpkinUID, ninjaSetupSrc)
		addWaveStrikerFillers(t, scn, player, 2)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 6000, scn.Match.GetPower(pumpkin, false))
	})

	t.Run("a small blocker cannot stop it once the ability is on", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		pumpkin := putCardInBattlezone(t, scn, player.Player, ninjaPumpkinUID, ninjaSetupSrc)
		addWaveStrikerFillers(t, scn, player, 2)
		blocker := putCardInBattlezone(t, scn, opponent.Player, waveStrikerSmallBlockerUID, ninjaSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		blockStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, pumpkin.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(opponent, blockStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action", "a 2000 power blocker is not offered the block")
		assert.Equal(t, match.BATTLEZONE, blocker.Zone)
	})
}
