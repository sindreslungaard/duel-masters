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
	minelordSkyterrorUID      = "0f74de48-a8da-4ff7-929a-50fa19fa1028"
	minelordSkyterror2000UID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	minelordSkyterror3000UID  = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur (3000)
	minelordSkyterror4000UID  = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	minelordSkyterrorSetupSrc = "minelord_skyterror_test_setup"
)

func TestMinelordSkyterror(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		minelord := putCardInBattlezone(t, scn, player.Player, minelordSkyterrorUID, minelordSkyterrorSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Minelord Skyterror", minelord.Name)
		assert.Equal(t, 3000, minelord.Power)
		assert.Equal(t, 4, minelord.ManaCost)
		assert.Equal(t, []string{civ.Fire}, minelord.Civs)
		assert.Equal(t, []string{civ.Fire}, minelord.ManaRequirement)
		assert.True(t, minelord.HasFamily(family.ArmoredWyvern))
		assert.True(t, minelord.HasCondition(cnd.SilentSkill))
	})

	t.Run("destroys every creature with power 3000 or less on both sides", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		minelord := putCardInBattlezone(t, scn, player.Player, minelordSkyterrorUID, minelordSkyterrorSetupSrc)
		minelord.Tapped = true

		ownSmall := putCardInBattlezone(t, scn, player.Player, minelordSkyterror2000UID, minelordSkyterrorSetupSrc)
		ownBig := putCardInBattlezone(t, scn, player.Player, minelordSkyterror4000UID, minelordSkyterrorSetupSrc)
		theirEdge := putCardInBattlezone(t, scn, opponent.Player, minelordSkyterror3000UID, minelordSkyterrorSetupSrc)
		theirBig := putCardInBattlezone(t, scn, opponent.Player, minelordSkyterror4000UID, minelordSkyterrorSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.GRAVEYARD, ownSmall.Zone)
		assert.Equal(t, match.GRAVEYARD, theirEdge.Zone, "exactly 3000 is within range")
		assert.Equal(t, match.BATTLEZONE, ownBig.Zone)
		assert.Equal(t, match.BATTLEZONE, theirBig.Zone)
		assert.Equal(t, match.GRAVEYARD, minelord.Zone, "at 3000 power it is caught by its own ability")
	})

	t.Run("declining destroys nothing", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		minelord := putCardInBattlezone(t, scn, player.Player, minelordSkyterrorUID, minelordSkyterrorSetupSrc)
		minelord.Tapped = true

		theirEdge := putCardInBattlezone(t, scn, opponent.Player, minelordSkyterror3000UID, minelordSkyterrorSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		declineSilentSkill(t, scn, player)

		assert.Equal(t, match.BATTLEZONE, theirEdge.Zone)
		assert.Equal(t, match.BATTLEZONE, minelord.Zone)
	})

	t.Run("it measures current power, not printed power", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		minelord := putCardInBattlezone(t, scn, player.Player, minelordSkyterrorUID, minelordSkyterrorSetupSrc)
		minelord.Tapped = true

		boosted := putCardInBattlezone(t, scn, opponent.Player, minelordSkyterror2000UID, minelordSkyterrorSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		// Applied after the untap step rebuilt conditions, so it is still in
		// place when the ability resolves at the start of the turn.
		boosted.AddUniqueSourceCondition(cnd.PowerAmplifier, 2000, minelordSkyterrorSetupSrc)
		require.Equal(t, 4000, scn.Match.GetPower(boosted, false))

		useSilentSkill(t, scn, player)

		assert.Equal(t, match.BATTLEZONE, boosted.Zone, "the bonus lifts it out of range")
		assert.Equal(t, match.GRAVEYARD, minelord.Zone)
	})
}
