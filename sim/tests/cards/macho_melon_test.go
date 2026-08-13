package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"testing"

	"github.com/stretchr/testify/assert"
)

const machoMelonCardUID = "fa987e39-2955-4074-bcf2-b7888ae27319"

func TestMachoMelon(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		melon := putCardInBattlezone(t, scn, player.Player, machoMelonCardUID, "macho_melon_test_setup")
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, melon, "Macho Melon", 1000, 2, []string{civ.Nature})
		assert.True(t, melon.HasFamily(family.WildVeggies))
		assert.True(t, melon.HasCondition(cnd.WaveStriker))
	})

	t.Run("its power attacker only applies while attacking and in numbers", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		melon := putCardInBattlezone(t, scn, player.Player, machoMelonCardUID, "macho_melon_test_setup")
		addWaveStrikerFillers(t, scn, player, 2)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 1000, scn.Match.GetPower(melon, false))
		assert.Equal(t, 4000, scn.Match.GetPower(melon, true))
	})
}
