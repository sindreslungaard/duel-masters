package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	clonedBladeUID      = "b760d56c-b6cc-4af3-aa97-fb3c04c232ff"
	clonedBladeSmallUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	clonedBladeBigUID   = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	clonedBladeSrc      = "cloned_blade_test_setup"
)

func TestClonedBlade(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(clonedBladeUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Cloned Blade", spell.Name)
		assert.Equal(t, 5, spell.ManaCost)
		assert.Equal(t, []string{civ.Fire}, spell.Civs)
		assert.Equal(t, []string{civ.Fire}, spell.ManaRequirement)
	})

	t.Run("with no copies buried it destroys exactly one", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		first := putCardInBattlezone(t, scn, opponent.Player, clonedBladeSmallUID, clonedBladeSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, clonedBladeSmallUID, clonedBladeSrc)

		castSpell(t, scn, player, clonedBladeUID)

		answerInTurn(t, scn, player, first.ID)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, first.Zone)
		assert.Equal(t, match.BATTLEZONE, second.Zone, "one copy, one target")
	})

	t.Run("each buried copy buys another target", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		first := putCardInBattlezone(t, scn, opponent.Player, clonedBladeSmallUID, clonedBladeSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, clonedBladeSmallUID, clonedBladeSrc)
		third := putCardInBattlezone(t, scn, opponent.Player, clonedBladeSmallUID, clonedBladeSrc)

		// One in each graveyard, so two extra targets on top of the first.
		_, err := player.Player.SpawnCard(clonedBladeUID, match.GRAVEYARD)
		require.NoError(t, err)
		_, err = opponent.Player.SpawnCard(clonedBladeUID, match.GRAVEYARD)
		require.NoError(t, err)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		castSpell(t, scn, player, clonedBladeUID)

		// The latest action rather than the first: paying for the spell opened
		// a prompt of its own, which castSpell has already answered.
		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, action.MinSelections, "the first target is mandatory")
		assert.Equal(t, 3, action.MaxSelections, "and two more are on offer")

		answerInTurn(t, scn, player, first.ID, second.ID, third.ID)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, first.Zone)
		assert.Equal(t, match.GRAVEYARD, second.Zone)
		assert.Equal(t, match.GRAVEYARD, third.Zone)
	})

	t.Run("the extra targets may be left alone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		first := putCardInBattlezone(t, scn, opponent.Player, clonedBladeSmallUID, clonedBladeSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, clonedBladeSmallUID, clonedBladeSrc)

		_, err := player.Player.SpawnCard(clonedBladeUID, match.GRAVEYARD)
		require.NoError(t, err)

		castSpell(t, scn, player, clonedBladeUID)

		answerInTurn(t, scn, player, first.ID)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, first.Zone)
		assert.Equal(t, match.BATTLEZONE, second.Zone, "\"you may choose another\" means it can be declined")
	})

	t.Run("anything above the power ceiling is safe", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		big := putCardInBattlezone(t, scn, opponent.Player, clonedBladeBigUID, clonedBladeSrc)
		small := putCardInBattlezone(t, scn, opponent.Player, clonedBladeSmallUID, clonedBladeSrc)

		castSpell(t, scn, player, clonedBladeUID)

		// Only one legal target and the choice is mandatory, so no prompt opens.
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, small.Zone)
		assert.Equal(t, match.BATTLEZONE, big.Zone, "4000 is above 3000")
	})

	t.Run("its caster's own creatures are never targets", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		mine := putCardInBattlezone(t, scn, player.Player, clonedBladeSmallUID, clonedBladeSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, clonedBladeSmallUID, clonedBladeSrc)

		castSpell(t, scn, player, clonedBladeUID)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, theirs.Zone)
		assert.Equal(t, match.BATTLEZONE, mine.Zone)
	})

	t.Run("an empty board asks nothing", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		castSpell(t, scn, player, clonedBladeUID)
		settleTurn(t, scn)
	})
}
