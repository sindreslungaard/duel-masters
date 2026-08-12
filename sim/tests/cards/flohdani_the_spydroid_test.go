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
	flohdaniTheSpydroidUID      = "476c377a-f91b-4c04-93ee-8c0c2ef27c5f"
	flohdaniTheSpydroidTargetID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	flohdaniTheSpydroidSetupSrc = "flohdani_the_spydroid_test_setup"
)

func TestFlohdaniTheSpydroid(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		flohdani := putCardInBattlezone(t, scn, player.Player, flohdaniTheSpydroidUID, flohdaniTheSpydroidSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Flohdani, the Spydroid", flohdani.Name)
		assert.Equal(t, 4000, flohdani.Power)
		assert.Equal(t, 4, flohdani.ManaCost)
		assert.Equal(t, []string{civ.Light}, flohdani.Civs)
		assert.Equal(t, []string{civ.Light}, flohdani.ManaRequirement)
		assert.True(t, flohdani.HasFamily(family.Soltrooper))
		assert.True(t, flohdani.HasCondition(cnd.SilentSkill))
	})

	t.Run("taps two of the opponent's creatures", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		flohdani := putCardInBattlezone(t, scn, player.Player, flohdaniTheSpydroidUID, flohdaniTheSpydroidSetupSrc)
		flohdani.Tapped = true

		first := putCardInBattlezone(t, scn, opponent.Player, flohdaniTheSpydroidTargetID, flohdaniTheSpydroidSetupSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, flohdaniTheSpydroidTargetID, flohdaniTheSpydroidSetupSrc)
		third := putCardInBattlezone(t, scn, opponent.Player, flohdaniTheSpydroidTargetID, flohdaniTheSpydroidSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		require.NoError(t, scn.SubmitAction(player, first.ID, second.ID))
		settleTurn(t, scn)

		assert.True(t, first.Tapped)
		assert.True(t, second.Tapped)
		assert.False(t, third.Tapped, "only the chosen creatures are tapped")
		assert.True(t, flohdani.Tapped)
	})

	t.Run("up to means one is allowed", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		flohdani := putCardInBattlezone(t, scn, player.Player, flohdaniTheSpydroidUID, flohdaniTheSpydroidSetupSrc)
		flohdani.Tapped = true

		first := putCardInBattlezone(t, scn, opponent.Player, flohdaniTheSpydroidTargetID, flohdaniTheSpydroidSetupSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, flohdaniTheSpydroidTargetID, flohdaniTheSpydroidSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		require.NoError(t, scn.SubmitAction(player, first.ID))
		settleTurn(t, scn)

		assert.True(t, first.Tapped)
		assert.False(t, second.Tapped)
	})

	t.Run("up to means none is allowed", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		flohdani := putCardInBattlezone(t, scn, player.Player, flohdaniTheSpydroidUID, flohdaniTheSpydroidSetupSrc)
		flohdani.Tapped = true

		target := putCardInBattlezone(t, scn, opponent.Player, flohdaniTheSpydroidTargetID, flohdaniTheSpydroidSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		require.NoError(t, scn.CancelAction(player))
		settleTurn(t, scn)

		assert.False(t, target.Tapped)
		assert.True(t, flohdani.Tapped, "the creature is kept tapped even if no target is chosen")
	})

	t.Run("an empty opposing battle zone opens no prompt", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		flohdani := putCardInBattlezone(t, scn, player.Player, flohdaniTheSpydroidUID, flohdaniTheSpydroidSetupSrc)
		flohdani.Tapped = true

		theirs, err := opponent.Player.Container(match.BATTLEZONE)
		require.NoError(t, err)
		require.Empty(t, theirs)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.True(t, flohdani.Tapped)
	})

	t.Run("it does not tap its controller's own creatures", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		flohdani := putCardInBattlezone(t, scn, player.Player, flohdaniTheSpydroidUID, flohdaniTheSpydroidSetupSrc)
		flohdani.Tapped = true

		own := putCardInBattlezone(t, scn, player.Player, flohdaniTheSpydroidTargetID, flohdaniTheSpydroidSetupSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, flohdaniTheSpydroidTargetID, flohdaniTheSpydroidSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		require.NoError(t, scn.SubmitAction(player, theirs.ID))
		settleTurn(t, scn)

		assert.True(t, theirs.Tapped)
		assert.False(t, own.Tapped)
	})
}
