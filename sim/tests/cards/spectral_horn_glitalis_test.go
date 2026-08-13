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
	spectralHornGlitalisUID      = "709b4470-2bad-453d-b0db-1cc218952403"
	spectralHornGlitalisSetupSrc = "spectral_horn_glitalis_test_setup"
)

func TestSpectralHornGlitalis(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, spectralHornGlitalisUID, spectralHornGlitalisSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Spectral Horn Glitalis", 4000, 3, []string{civ.Nature, civ.Light})
		assert.True(t, card.HasFamily(family.HornedBeast))
		assert.True(t, card.IsMulticolored())
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(spectralHornGlitalisUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, spectralHornGlitalisSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})
}
