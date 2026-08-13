package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	electroExplorerSyrionUID      = "83f689d5-5653-4306-84c7-c43d3dd390a7"
	electroExplorerSyrionSetupSrc = "electro_explorer_syrion_test_setup"
)

func TestElectroExplorerSyrion(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, electroExplorerSyrionUID, electroExplorerSyrionSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Electro Explorer Syrion", 4000, 3, []string{civ.Light, civ.Water})
		assert.True(t, card.HasFamily(family.Gladiator))
		assert.True(t, card.IsMulticolored())
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(electroExplorerSyrionUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, electroExplorerSyrionSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})
}
