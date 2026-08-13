package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	necrodragonJagraveenUID      = "8558c3cc-b2ef-4d54-aec8-a4b5bde0b904"
	necrodragonJagraveenAttacker = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (fire, cost 2)
	necrodragonJagraveenSetupSrc = "necrodragon_jagraveen_test_setup"
)

func TestNecrodragonJagraveen(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, necrodragonJagraveenUID, necrodragonJagraveenSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Necrodragon Jagraveen", 6000, 6, []string{civ.Darkness})
		assert.True(t, card.HasFamily(family.ZombieDragon))
		assert.True(t, card.HasCondition(cnd.Blocker))
		assert.True(t, card.HasCondition(cnd.DoubleBreaker))
		assert.False(t, card.HasCondition(cnd.DestroyAfterBattle), "only blocking costs it its life")
	})

	t.Run("blocking destroys it even when it wins", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		jagraveen := putCardInBattlezone(t, scn, player.Player, necrodragonJagraveenUID, necrodragonJagraveenSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, necrodragonJagraveenAttacker, necrodragonJagraveenSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		blockStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(opponent, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(opponent, shields[0].ID))

		_, err = scn.WaitForAction(player, blockStart)
		require.NoError(t, err)
		answerInTurn(t, scn, player, jagraveen.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, attacker.Zone, "6000 beats 2000")
		assert.Equal(t, match.GRAVEYARD, jagraveen.Zone, "and blocking costs it its life")
	})

	t.Run("attacking does not destroy it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		jagraveen := putCardInBattlezone(t, scn, player.Player, necrodragonJagraveenUID, necrodragonJagraveenSetupSrc)
		victim := putCardInBattlezone(t, scn, opponent.Player, necrodragonJagraveenAttacker, necrodragonJagraveenSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		victim.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(player, jagraveen.ID, victim.ID))

		assert.Equal(t, match.GRAVEYARD, victim.Zone)
		assert.Equal(t, match.BATTLEZONE, jagraveen.Zone, "the printed trigger is blocking, not battling")
	})
}
