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
	seaMutantDormelUID      = "854a1e8c-cebb-4d19-b22c-d40fa3761bf1"
	seaMutantDormelSetupSrc = "sea_mutant_dormel_test_setup"
)

func TestSeaMutantDormel(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, seaMutantDormelUID, seaMutantDormelSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Sea Mutant Dormel", 4000, 3, []string{civ.Water, civ.Darkness})
		assert.True(t, card.HasFamily(family.Merfolk))
		assert.True(t, card.IsMulticolored())
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(seaMutantDormelUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, seaMutantDormelSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})
}
