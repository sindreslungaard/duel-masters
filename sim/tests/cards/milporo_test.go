package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	milporoUID      = "460fc2eb-c7cd-42d5-9bed-a98de4f59026"
	milporoSetupSrc = "milporo_test_setup"
)

func TestMilporo(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		milporo := putCardInBattlezone(t, scn, player.Player, milporoUID, milporoSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Milporo", milporo.Name)
		assert.Equal(t, 3000, milporo.Power)
		assert.Equal(t, 4, milporo.ManaCost)
		assert.Equal(t, []string{civ.Water}, milporo.Civs)
		assert.Equal(t, []string{civ.Water}, milporo.ManaRequirement)
		assert.True(t, milporo.HasFamily(family.CyberLord))
		assert.True(t, milporo.HasCondition(cnd.SilentSkill))
	})

	t.Run("draws a card on top of the draw step", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		milporo := putCardInBattlezone(t, scn, player.Player, milporoUID, milporoSetupSrc)
		milporo.Tapped = true

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore)+2)
		assert.True(t, milporo.Tapped, "it stays tapped for the turn")
	})

	t.Run("declining draws only for the draw step", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		milporo := putCardInBattlezone(t, scn, player.Player, milporoUID, milporoSetupSrc)
		milporo.Tapped = true

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		declineSilentSkill(t, scn, player)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore)+1)
		assert.False(t, milporo.Tapped, "it untaps and can attack instead")
	})
}
