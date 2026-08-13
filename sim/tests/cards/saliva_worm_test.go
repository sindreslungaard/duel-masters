package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	salivaWormUID   = "2ea5870b-6f0f-4980-89de-c5c2d50c798f"
	salivaWormSetup = "saliva_worm_test_setup"
)

func TestSalivaWorm(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		worm := putCardInBattlezone(t, scn, player.Player, salivaWormUID, salivaWormSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, worm, "Saliva Worm", 2000, 3, []string{civ.Darkness})
		assert.True(t, worm.HasFamily(family.ParasiteWorm))
		assert.True(t, worm.HasCondition(cnd.WaveStriker))
		assert.False(t, worm.HasCondition(cnd.Stealth), "no stealth without the count")
	})

	t.Run("with the count it grows and gains darkness stealth", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		worm := putCardInBattlezone(t, scn, player.Player, salivaWormUID, salivaWormSetup)
		addWaveStrikerFillers(t, scn, player, 2)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 6000, scn.Match.GetPower(worm, false))
		assert.True(t, worm.HasCondition(cnd.Stealth))
	})
}
