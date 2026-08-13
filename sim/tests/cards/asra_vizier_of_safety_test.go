package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	asraVizierOfSafetyUID = "9856078c-0319-4867-a190-427e02c043bb"
	asraSetupSrc          = "asra_vizier_of_safety_test_setup"
)

func TestAsraVizierOfSafety(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		asra := putCardInBattlezone(t, scn, player.Player, asraVizierOfSafetyUID, asraSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, asra, "Asra, Vizier of Safety", 2000, 3, []string{civ.Light})
		assert.True(t, asra.HasFamily(family.Initiate))
		assert.True(t, asra.HasCondition(cnd.WaveStriker))
	})

	t.Run("without the count it is neither bigger nor a blocker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		asra := putCardInBattlezone(t, scn, player.Player, asraVizierOfSafetyUID, asraSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 2000, scn.Match.GetPower(asra, false))
		assert.False(t, asra.HasCondition(cnd.Blocker))
	})

	t.Run("with the count it gains both halves", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		asra := putCardInBattlezone(t, scn, player.Player, asraVizierOfSafetyUID, asraSetupSrc)
		addWaveStrikerFillers(t, scn, player, 2)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 6000, scn.Match.GetPower(asra, false), "+4000")
		assert.True(t, asra.HasCondition(cnd.Blocker))
	})
}
