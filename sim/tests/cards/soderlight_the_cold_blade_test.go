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
	soderlightTheColdBladeUID      = "868dd8dd-7777-4bc2-94d4-7d7ccaf8f999"
	soderlightTheColdBladePlainUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	soderlightTheColdBladeSetupSrc = "soderlight_the_cold_blade_test_setup"
)

func TestSoderlightTheColdBlade(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		soderlight := putCardInBattlezone(t, scn, player.Player, soderlightTheColdBladeUID, soderlightTheColdBladeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Soderlight, the Cold Blade", soderlight.Name)
		assert.Equal(t, 4000, soderlight.Power)
		assert.Equal(t, 4, soderlight.ManaCost)
		assert.Equal(t, []string{civ.Water, civ.Darkness}, soderlight.Civs)
		assert.Equal(t, []string{civ.Water, civ.Darkness}, soderlight.ManaRequirement)
		assert.True(t, soderlight.IsMulticolored())
		assert.True(t, soderlight.HasFamily(family.SpiritQuartz))
		assert.True(t, soderlight.HasCondition(cnd.CantBeBlocked))
		assert.True(t, soderlight.HasCondition(cnd.SilentSkill))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupSilentSkillTest(t)

		card, err := player.Player.SpawnCard(soderlightTheColdBladeUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, soderlightTheColdBladeSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.MANAZONE, moved.Zone)

		assert.True(t, moved.Tapped, "a multicolored card is charged tapped")
	})

	t.Run("the opponent chooses which of their creatures dies", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		soderlight := putCardInBattlezone(t, scn, player.Player, soderlightTheColdBladeUID, soderlightTheColdBladeSetupSrc)
		soderlight.Tapped = true

		first := putCardInBattlezone(t, scn, opponent.Player, soderlightTheColdBladePlainUID, soderlightTheColdBladeSetupSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, soderlightTheColdBladePlainUID, soderlightTheColdBladeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		// The prompt belongs to the opponent, not to the card's controller.
		require.NoError(t, scn.SubmitAction(opponent, second.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, second.Zone)
		assert.Equal(t, match.BATTLEZONE, first.Zone)
		assert.True(t, soderlight.Tapped)
	})

	t.Run("a single creature is destroyed without asking", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		soderlight := putCardInBattlezone(t, scn, player.Player, soderlightTheColdBladeUID, soderlightTheColdBladeSetupSrc)
		soderlight.Tapped = true

		only := putCardInBattlezone(t, scn, opponent.Player, soderlightTheColdBladePlainUID, soderlightTheColdBladeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.GRAVEYARD, only.Zone)
	})

	t.Run("an opponent with no creatures loses nothing", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		soderlight := putCardInBattlezone(t, scn, player.Player, soderlightTheColdBladeUID, soderlightTheColdBladeSetupSrc)
		soderlight.Tapped = true

		own := putCardInBattlezone(t, scn, player.Player, soderlightTheColdBladePlainUID, soderlightTheColdBladeSetupSrc)

		theirs, err := opponent.Player.Container(match.BATTLEZONE)
		require.NoError(t, err)
		require.Empty(t, theirs)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.BATTLEZONE, own.Zone, "the controller's own creatures are never eligible")
		assert.Equal(t, match.BATTLEZONE, soderlight.Zone)
		assert.True(t, soderlight.Tapped)
	})

	t.Run("declining destroys nothing", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		soderlight := putCardInBattlezone(t, scn, player.Player, soderlightTheColdBladeUID, soderlightTheColdBladeSetupSrc)
		soderlight.Tapped = true

		theirs := putCardInBattlezone(t, scn, opponent.Player, soderlightTheColdBladePlainUID, soderlightTheColdBladeSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		declineSilentSkill(t, scn, player)

		assert.Equal(t, match.BATTLEZONE, theirs.Zone)
		assert.False(t, soderlight.Tapped)
	})
}
