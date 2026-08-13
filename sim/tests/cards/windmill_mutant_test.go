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
	windmillMutantUID      = "88435ced-6e13-4084-ba0e-21cac66808e9"
	windmillMutantSetupSrc = "windmill_mutant_test_setup"
)

func TestWindmillMutant(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, windmillMutantUID, windmillMutantSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Windmill Mutant", 2000, 3, []string{civ.Darkness})
		assert.True(t, card.HasFamily(family.Hedrian))
	})

	t.Run("attacking costs the opponent a card at random", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		windmill := putCardInBattlezone(t, scn, player.Player, windmillMutantUID, windmillMutantSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, windmill.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())

		// One card lost to the discard, one gained from the broken shield.
		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore))
		assert.Equal(t, match.HAND, shields[0].Zone)
	})
}
