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
	lockdownLizardUID = "5369887a-6aaa-489a-8aa2-543e74162832"
	lockdownTapperUID = "a808b98c-2de7-412b-970c-a3b925bf43c2" // Deklowaz, the Terminator (tap ability)
	lockdownVictimUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	lockdownSetupSrc  = "lockdown_lizard_test_setup"
)

func TestLockdownLizard(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lizard := putCardInBattlezone(t, scn, player.Player, lockdownLizardUID, lockdownSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, lizard, "Lockdown Lizard", 3000, 4, []string{civ.Fire})
		assert.True(t, lizard.HasFamily(family.MeltWarrior))
	})

	t.Run("the opponent cannot use a tap ability", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, lockdownLizardUID, lockdownSetupSrc)
		tapper := putCardInBattlezone(t, scn, opponent.Player, lockdownTapperUID, lockdownSetupSrc)
		victim := putCardInBattlezone(t, scn, player.Player, lockdownVictimUID, lockdownSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		require.True(t, tapper.HasCondition(cnd.TapAbility))

		require.Error(t, scn.ActionUseTapAbility(opponent, tapper.ID), "the ability should be refused")

		assert.False(t, tapper.Tapped, "a refused ability does not tap the creature")
		assert.Equal(t, match.BATTLEZONE, victim.Zone, "and its effect never resolved")
	})

	t.Run("its own controller is locked down too", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, lockdownLizardUID, lockdownSetupSrc)
		tapper := putCardInBattlezone(t, scn, player.Player, lockdownTapperUID, lockdownSetupSrc)
		victim := putCardInBattlezone(t, scn, opponent.Player, lockdownVictimUID, lockdownSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		require.Error(t, scn.ActionUseTapAbility(player, tapper.ID), "\"players\" is both of them")

		assert.False(t, tapper.Tapped)
		assert.Equal(t, match.BATTLEZONE, victim.Zone)
	})

	t.Run("the lock lifts when it leaves the battle zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lizard := putCardInBattlezone(t, scn, player.Player, lockdownLizardUID, lockdownSetupSrc)
		tapper := putCardInBattlezone(t, scn, player.Player, lockdownTapperUID, lockdownSetupSrc)
		victim := putCardInBattlezone(t, scn, opponent.Player, lockdownVictimUID, lockdownSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.Error(t, scn.ActionUseTapAbility(player, tapper.ID))

		_, err := player.Player.MoveCard(lizard.ID, match.BATTLEZONE, match.GRAVEYARD, lockdownSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		require.NoError(t, scn.ActionUseTapAbility(player, tapper.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, tapper.Tapped)
		assert.Equal(t, match.GRAVEYARD, victim.Zone, "the ability resolved once the lock was gone")
	})

	t.Run("attacking is untouched", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, lockdownLizardUID, lockdownSetupSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, lockdownTapperUID, lockdownSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err, "only tap abilities are locked down")
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))
	})
}
