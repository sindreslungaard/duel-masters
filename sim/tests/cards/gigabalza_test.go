package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	gigabalzaUID      = "2f6d65f6-6ec7-4ee5-8c8f-732938b8eeb6"
	gigabalzaSetupSrc = "gigabalza_test_setup"
)

func TestGigabalza(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, gigabalzaUID, gigabalzaSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Gigabalza", 1000, 4, []string{civ.Darkness})
		assert.True(t, card.HasFamily(family.Chimera))
		assert.True(t, card.HasCondition(cnd.ShieldTrigger))
	})

	t.Run("arriving costs the opponent a card at random", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		gigabalza := spawnForLater(t, player, gigabalzaUID)
		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		putIntoPlay(t, scn, player, gigabalza)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore)-1)
	})
}
