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
	gigappiPontoUID      = "bb3a1ebc-2b0f-4b7e-8c44-aa3bf1e66931"
	gigappiPontoSetupSrc = "gigappi_ponto_test_setup"
)

func TestGigappiPonto(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, gigappiPontoUID, gigappiPontoSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Gigappi Ponto", 4000, 3, []string{civ.Darkness, civ.Fire})
		assert.True(t, card.HasFamily(family.Chimera))
		assert.True(t, card.IsMulticolored())
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(gigappiPontoUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, gigappiPontoSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})
}
