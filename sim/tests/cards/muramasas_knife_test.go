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
	muramasasKnifeUID       = "7f2e9361-2399-4ec1-a2b2-af1c88288833"
	muramasasKnifeSetupSrc  = "muramasasKnife_test_setup"
	muramasasKnifeVictimUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (fire, cost 2)
)

func TestMuramasasKnife(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, muramasasKnifeUID, muramasasKnifeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Muramasa's Knife", 2000, 3, []string{civ.Fire})
		assert.True(t, card.HasFamily(family.Xenoparts))
		assert.True(t, card.HasCondition(cnd.AttackUntapped))
	})

	t.Run("it can attack an untapped creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, muramasasKnifeUID, muramasasKnifeSetupSrc)
		victim := putCardInBattlezone(t, scn, opponent.Player, muramasasKnifeVictimUID, muramasasKnifeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.False(t, victim.Tapped, "normally untouchable")

		require.NoError(t, scn.ActionAttackCreature(player, card.ID, victim.ID))

		// Both are 2000, so they trade.
		assert.Equal(t, match.GRAVEYARD, victim.Zone)
	})
}
