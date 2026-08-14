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
	necrodragonZalvaUID      = "c17bd6a6-c5ae-4235-b0b1-b5e8321d3a06"
	necrodragonZalvaSetupSrc = "necrodragon_zalva_test_setup"
)

func TestNecrodragonZalva(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, necrodragonZalvaUID, necrodragonZalvaSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Necrodragon Zalva", 5000, 4, []string{civ.Darkness})
		assert.True(t, card.HasFamily(family.ZombieDragon))
	})

	t.Run("arriving lets the opponent draw", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		zalva := spawnForLater(t, player, necrodragonZalvaUID)
		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		putIntoPlay(t, scn, player, zalva)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore)+1)
	})
}
