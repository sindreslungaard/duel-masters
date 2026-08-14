package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	valkyerStarstormUID      = "51ce0197-7d12-498a-a874-17e60a0d3f21"
	valkyerStarstormSetupSrc = "valkyerStarstorm_test_setup"
)

func TestValkyerStarstormElemental(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, valkyerStarstormUID, valkyerStarstormSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Valkyer, Starstorm Elemental", 7000, 5, []string{civ.Light})
		assert.True(t, card.HasFamily(family.AngelCommand))
		assert.True(t, card.HasCondition(cnd.Blocker))
		assert.True(t, card.HasCondition(cnd.CantAttackPlayers))
	})

	t.Run("it cannot attack the player", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, valkyerStarstormUID, valkyerStarstormSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		_, err := scn.ActionAttackPlayer(player, card.ID)
		assert.Error(t, err)
		assert.False(t, card.Tapped)
	})
}
