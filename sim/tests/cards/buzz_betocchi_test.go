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
	buzzBetocchiUID      = "124aff9d-aaca-4c88-b40b-2529127a1214"
	buzzBetocchiSetupSrc = "buzz_betocchi_test_setup"
)

func TestBuzzBetocchi(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, buzzBetocchiUID, buzzBetocchiSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Buzz Betocchi", 4000, 3, []string{civ.Fire, civ.Nature})
		assert.True(t, card.HasFamily(family.FireBird))
		assert.True(t, card.IsMulticolored())
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(buzzBetocchiUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, buzzBetocchiSetupSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})
}
