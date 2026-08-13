package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	peppiPepperUID      = "ce3863f2-5810-40da-9c55-297affdd5787"
	peppiPepperSetupSrc = "peppiPepper_test_setup"
)

func TestPeppiPepper(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, peppiPepperUID, peppiPepperSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Peppi Pepper", 2000, 3, []string{civ.Fire})
		assert.True(t, card.HasFamily(family.FireBird))
		assert.True(t, card.HasCondition(cnd.PowerAttacker))
	})

	t.Run("power attacker only counts while attacking", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, peppiPepperUID, peppiPepperSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 2000, scn.Match.GetPower(card, false))
		assert.Equal(t, 5000, scn.Match.GetPower(card, true))
	})
}
