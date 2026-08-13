package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baraidTheExplorerUID = "7b3480ab-57fb-413e-bd93-7ffdfeb2d73f"
	baraidLightAllyUID   = "6bfebd3d-64ff-4321-bb0b-e994ca8f811e" // Telitol, the Explorer (light)
	baraidSetupSrc       = "baraid_the_explorer_test_setup"
)

func TestBaraidTheExplorer(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		baraid := putCardInBattlezone(t, scn, player.Player, baraidTheExplorerUID, baraidSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, baraid, "Baraid, the Explorer", 5000, 5, []string{civ.Light})
		assert.True(t, baraid.HasFamily(family.Gladiator))
		assert.True(t, baraid.HasCondition(cnd.SilentSkill))
	})

	t.Run("its silent skill makes the light creatures unblockable", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		baraid := putCardInBattlezone(t, scn, player.Player, baraidTheExplorerUID, baraidSetupSrc)
		baraid.Tapped = true

		lightAlly := putCardInBattlezone(t, scn, player.Player, baraidLightAllyUID, baraidSetupSrc)
		fireAlly := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, baraidSetupSrc)
		theirLight := putCardInBattlezone(t, scn, opponent.Player, baraidLightAllyUID, baraidSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.True(t, lightAlly.HasCondition(cnd.CantBeBlocked))
		assert.True(t, baraid.HasCondition(cnd.CantBeBlocked), "it is a light creature of its own")
		assert.False(t, fireAlly.HasCondition(cnd.CantBeBlocked), "only light creatures")
		assert.False(t, theirLight.HasCondition(cnd.CantBeBlocked), "only its controller's creatures")
	})

	t.Run("the grant is gone by the next turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		baraid := putCardInBattlezone(t, scn, player.Player, baraidTheExplorerUID, baraidSetupSrc)
		baraid.Tapped = true

		lightAlly := putCardInBattlezone(t, scn, player.Player, baraidLightAllyUID, baraidSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)
		require.True(t, lightAlly.HasCondition(cnd.CantBeBlocked))

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		// It is granted again each turn, but the previous turn's copy did not
		// survive: the condition carries only this turn's source.
		assert.True(t, lightAlly.HasCondition(cnd.CantBeBlocked))
	})

	t.Run("declining grants nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		baraid := putCardInBattlezone(t, scn, player.Player, baraidTheExplorerUID, baraidSetupSrc)
		baraid.Tapped = true

		lightAlly := putCardInBattlezone(t, scn, player.Player, baraidLightAllyUID, baraidSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		declineSilentSkill(t, scn, player)

		assert.False(t, lightAlly.HasCondition(cnd.CantBeBlocked))
	})
}
