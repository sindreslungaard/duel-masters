package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	merleeTheOracleUID = "9ec241a3-57e9-4054-8680-aff2c1a7b45b"
	merleeSetupSrc     = "merlee_the_oracle_test_setup"
)

func TestMerleeTheOracle(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		merlee := putCardInBattlezone(t, scn, player.Player, merleeTheOracleUID, merleeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, merlee, "Merlee, the Oracle", 1500, 2, []string{civ.Light})
		assert.True(t, merlee.HasFamily(family.LightBringer))
		assert.True(t, merlee.HasCondition(cnd.WaveStriker))
	})

	t.Run("it lifts every creature its controller has", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		merlee := putCardInBattlezone(t, scn, player.Player, merleeTheOracleUID, merleeSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, merleeSetupSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, merleeSetupSrc)
		addWaveStrikerFillers(t, scn, player, 2)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 3000, scn.Match.GetPower(ally, false))
		assert.Equal(t, 2500, scn.Match.GetPower(merlee, false), "\"each of your creatures\" includes itself")
		assert.Equal(t, 2000, scn.Match.GetPower(theirs, false), "not the opponent's")
	})

	t.Run("without the count nothing is lifted", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, merleeTheOracleUID, merleeSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, merleeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 2000, scn.Match.GetPower(ally, false))
	})
}
