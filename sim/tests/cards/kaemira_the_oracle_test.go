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
	kaemiraTheOracleUID      = "afeecd61-4c89-4ab9-9c25-297551fe3624"
	kaemiraTheOracleSeedUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	kaemiraTheOracleSetupSrc = "kaemira_the_oracle_test_setup"
)

func TestKaemiraTheOracle(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		kaemira := putCardInBattlezone(t, scn, player.Player, kaemiraTheOracleUID, kaemiraTheOracleSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Kaemira, the Oracle", kaemira.Name)
		assert.Equal(t, 1000, kaemira.Power)
		assert.Equal(t, 4, kaemira.ManaCost)
		assert.Equal(t, []string{civ.Light}, kaemira.Civs)
		assert.Equal(t, []string{civ.Light}, kaemira.ManaRequirement)
		assert.True(t, kaemira.HasFamily(family.LightBringer))
		assert.True(t, kaemira.HasCondition(cnd.SilentSkill))
	})

	t.Run("adds the top card of the deck to the shields", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		kaemira := putCardInBattlezone(t, scn, player.Player, kaemiraTheOracleUID, kaemiraTheOracleSetupSrc)
		kaemira.Tapped = true

		player.Player.DestroyDeck()
		for range 4 {
			_, err := player.Player.SpawnCard(kaemiraTheOracleSeedUID, match.DECK)
			require.NoError(t, err)
		}

		topBefore := player.Player.PeekDeck(1)
		require.Len(t, topBefore, 1)

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore)+1)
		assert.Equal(t, match.SHIELDZONE, topBefore[0].Zone)
		assert.True(t, kaemira.Tapped)
	})

	t.Run("declining leaves the shields alone", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		kaemira := putCardInBattlezone(t, scn, player.Player, kaemiraTheOracleUID, kaemiraTheOracleSetupSrc)
		kaemira.Tapped = true

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		declineSilentSkill(t, scn, player)

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore))
		assert.False(t, kaemira.Tapped)
	})
}
