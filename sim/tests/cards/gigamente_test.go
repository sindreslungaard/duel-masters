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
	gigamenteUID          = "e7b18dba-e6e6-426a-9e57-c4b1088d79d6"
	gigamenteCreatureUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	gigamenteSpellUID     = "5883180e-d88c-4f24-b17c-f5a837420147" // Terror Pit
	gigamenteTestSetupSrc = "gigamente_test_setup"
)

func TestGigamente(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		gigamente := putCardInBattlezone(t, scn, player.Player, gigamenteUID, gigamenteTestSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Gigamente", gigamente.Name)
		assert.Equal(t, 3000, gigamente.Power)
		assert.Equal(t, 4, gigamente.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, gigamente.Civs)
		assert.Equal(t, []string{civ.Darkness}, gigamente.ManaRequirement)
		assert.True(t, gigamente.HasFamily(family.Chimera))
		assert.True(t, gigamente.HasCondition(cnd.SilentSkill))
	})

	t.Run("returns the chosen creature from the graveyard", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		gigamente := putCardInBattlezone(t, scn, player.Player, gigamenteUID, gigamenteTestSetupSrc)
		gigamente.Tapped = true

		first, err := player.Player.SpawnCard(gigamenteCreatureUID, match.GRAVEYARD)
		require.NoError(t, err)
		second, err := player.Player.SpawnCard(gigamenteCreatureUID, match.GRAVEYARD)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		// Two candidates, so the player is asked which one comes back.
		require.NoError(t, scn.SubmitAction(player, second.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, second.Zone)
		assert.Equal(t, match.GRAVEYARD, first.Zone, "only the chosen creature returns")
	})

	t.Run("a single creature returns without asking", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		gigamente := putCardInBattlezone(t, scn, player.Player, gigamenteUID, gigamenteTestSetupSrc)
		gigamente.Tapped = true

		only, err := player.Player.SpawnCard(gigamenteCreatureUID, match.GRAVEYARD)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		// The return is mandatory once the ability is used, so with one legal
		// choice the engine takes it rather than opening a pointless prompt.
		assert.Equal(t, match.HAND, only.Zone)
	})

	t.Run("a graveyard without creatures returns nothing", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		gigamente := putCardInBattlezone(t, scn, player.Player, gigamenteUID, gigamenteTestSetupSrc)
		gigamente.Tapped = true

		spell, err := player.Player.SpawnCard(gigamenteSpellUID, match.GRAVEYARD)
		require.NoError(t, err)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.GRAVEYARD, spell.Zone, "a spell is not a creature")

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore)+1, "only the draw step")
		assert.True(t, gigamente.Tapped)
	})

	t.Run("an empty graveyard returns nothing", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		gigamente := putCardInBattlezone(t, scn, player.Player, gigamenteUID, gigamenteTestSetupSrc)
		gigamente.Tapped = true

		graveyard, err := player.Player.Container(match.GRAVEYARD)
		require.NoError(t, err)
		require.Empty(t, graveyard)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore)+1, "only the draw step")
		assert.True(t, gigamente.Tapped)
	})

	t.Run("it cannot pull a creature out of the opponent's graveyard", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		gigamente := putCardInBattlezone(t, scn, player.Player, gigamenteUID, gigamenteTestSetupSrc)
		gigamente.Tapped = true

		theirs, err := opponent.Player.SpawnCard(gigamenteCreatureUID, match.GRAVEYARD)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.GRAVEYARD, theirs.Zone)
	})
}
