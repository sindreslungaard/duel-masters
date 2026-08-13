package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	hypersprintUzesolUID      = "2622d460-6424-4410-9662-f0f4cd9c08b0"
	hypersprintUzesolSetupSrc = "hypersprintUzesol_test_setup"
)

func TestHypersprintWariorUzesol(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, hypersprintUzesolUID, hypersprintUzesolSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Hypersprint Warior Uzesol", 1000, 4, []string{civ.Fire})
		assert.True(t, card.HasFamily(family.Armorloid))
		assert.True(t, card.HasCondition(cnd.SpeedAttacker))
		assert.True(t, card.HasCondition(cnd.PowerAttacker))
	})

	t.Run("power attacker only counts while attacking", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, hypersprintUzesolUID, hypersprintUzesolSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 1000, scn.Match.GetPower(card, false))
		assert.Equal(t, 5000, scn.Match.GetPower(card, true))
	})
}
