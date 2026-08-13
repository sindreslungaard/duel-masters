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
	radioactiveHornUID      = "aec7a0eb-a16a-4a8a-bd18-e9adb8970432"
	radioactiveHornSetupSrc = "radioactiveHorn_test_setup"
)

func TestRadioactiveHornTheStrange(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, radioactiveHornUID, radioactiveHornSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Radioactive Horn, the Strange", 1000, 3, []string{civ.Nature})
		assert.True(t, card.HasFamily(family.HornedBeast))
		assert.True(t, card.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("it breaks two shields", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, radioactiveHornUID, radioactiveHornSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(shields), 2)

		_, err = scn.ActionAttackPlayer(player, card.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID, shields[1].ID))

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, len(shields)-2)
	})
}
