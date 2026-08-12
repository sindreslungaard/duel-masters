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
	vorgsEngineUID      = "1b786e62-ffbc-4694-a8cd-8dd48f8e18fd"
	vorgsEngineWeakUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	vorgsEngineToughUID = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur (3000)
	vorgsEngineSetupSrc = "vorgs_engine_test_setup"
)

func TestVorgsEngine(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		engine := putCardInBattlezone(t, scn, player.Player, vorgsEngineUID, vorgsEngineSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Vorg's Engine", engine.Name)
		assert.Equal(t, 2000, engine.Power)
		assert.Equal(t, 2, engine.ManaCost)
		assert.Equal(t, []string{civ.Fire}, engine.Civs)
		assert.Equal(t, []string{civ.Fire}, engine.ManaRequirement)
		assert.True(t, engine.HasFamily(family.Xenoparts))
		assert.True(t, engine.HasCondition(cnd.SilentSkill))
	})

	t.Run("destroys every creature with power 2000 or less, its own side included", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		engine := putCardInBattlezone(t, scn, player.Player, vorgsEngineUID, vorgsEngineSetupSrc)
		engine.Tapped = true

		ownWeak := putCardInBattlezone(t, scn, player.Player, vorgsEngineWeakUID, vorgsEngineSetupSrc)
		theirWeak := putCardInBattlezone(t, scn, opponent.Player, vorgsEngineWeakUID, vorgsEngineSetupSrc)
		theirTough := putCardInBattlezone(t, scn, opponent.Player, vorgsEngineToughUID, vorgsEngineSetupSrc)
		require.Greater(t, theirTough.Power, 2000, "the survivor has to be out of range")

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.GRAVEYARD, ownWeak.Zone)
		assert.Equal(t, match.GRAVEYARD, theirWeak.Zone)
		assert.Equal(t, match.BATTLEZONE, theirTough.Zone, "power above 2000 survives")
		assert.Equal(t, match.GRAVEYARD, engine.Zone, "at 2000 power it destroys itself too")
	})

	t.Run("declining destroys nothing", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		engine := putCardInBattlezone(t, scn, player.Player, vorgsEngineUID, vorgsEngineSetupSrc)
		engine.Tapped = true

		theirWeak := putCardInBattlezone(t, scn, opponent.Player, vorgsEngineWeakUID, vorgsEngineSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		declineSilentSkill(t, scn, player)

		assert.Equal(t, match.BATTLEZONE, theirWeak.Zone)
		assert.Equal(t, match.BATTLEZONE, engine.Zone)
		assert.False(t, engine.Tapped)
	})

	t.Run("an empty battle zone opposite is not a problem", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		engine := putCardInBattlezone(t, scn, player.Player, vorgsEngineUID, vorgsEngineSetupSrc)
		engine.Tapped = true

		theirs, err := opponent.Player.Container(match.BATTLEZONE)
		require.NoError(t, err)
		require.Empty(t, theirs)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.GRAVEYARD, engine.Zone)
	})
}
