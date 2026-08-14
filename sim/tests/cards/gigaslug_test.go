package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	gigaslugUID       = "a1e3713b-21f2-4fff-adf0-1bdabf4292d1"
	gigaslugSetupSrc  = "gigaslug_test_setup"
	gigaslugVictimUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (fire, cost 2)
)

func TestGigaslug(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, gigaslugUID, gigaslugSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Gigaslug", 1000, 3, []string{civ.Darkness})
		assert.True(t, card.HasFamily(family.Chimera))
		assert.True(t, card.HasCondition(cnd.Blocker))
		assert.True(t, card.HasCondition(cnd.Slayer))
		assert.True(t, card.HasCondition(cnd.CantAttackCreatures))
		assert.True(t, card.HasCondition(cnd.CantAttackPlayers))
	})

	t.Run("it cannot attack at all", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, gigaslugUID, gigaslugSetupSrc)
		victim := putCardInBattlezone(t, scn, opponent.Player, gigaslugVictimUID, gigaslugSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		victim.Tapped = true

		_, err := scn.ActionAttackPlayer(player, card.ID)
		require.Error(t, err)
		require.Error(t, scn.ActionAttackCreature(player, card.ID, victim.ID))

		assert.False(t, card.Tapped)
	})
}
