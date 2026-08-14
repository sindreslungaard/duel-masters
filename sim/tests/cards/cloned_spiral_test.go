package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	clonedSpiralUID         = "55367300-8226-4abf-bdce-061549630013"
	clonedSpiralCreatureUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	clonedSpiralSrc         = "cloned_spiral_test_setup"
)

func TestClonedSpiral(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(clonedSpiralUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Cloned Spiral", spell.Name)
		assert.Equal(t, 4, spell.ManaCost)
		assert.Equal(t, []string{civ.Water}, spell.Civs)
		assert.Equal(t, []string{civ.Water}, spell.ManaRequirement)
	})

	t.Run("with no copies buried it returns exactly one", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		first := putCardInBattlezone(t, scn, opponent.Player, clonedSpiralCreatureUID, clonedSpiralSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, clonedSpiralCreatureUID, clonedSpiralSrc)

		castSpell(t, scn, player, clonedSpiralUID)

		answerInTurn(t, scn, player, first.ID)
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, first.Zone)
		assert.Equal(t, match.BATTLEZONE, second.Zone)
	})

	t.Run("either side of the board is fair game", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		mine := putCardInBattlezone(t, scn, player.Player, clonedSpiralCreatureUID, clonedSpiralSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, clonedSpiralCreatureUID, clonedSpiralSrc)

		_, err := player.Player.SpawnCard(clonedSpiralUID, match.GRAVEYARD)
		require.NoError(t, err)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		castSpell(t, scn, player, clonedSpiralUID)

		action, err := scn.WaitForMultipartAction(player, promptStart)
		require.NoError(t, err)
		assert.Contains(t, action.Cards, "Your creatures")
		assert.Contains(t, action.Cards, "Your opponent's creatures")

		answerInTurn(t, scn, player, mine.ID, theirs.ID)
		settleTurn(t, scn)

		// Each goes to its own owner's hand, not to the caster's.
		assert.Equal(t, match.HAND, mine.Zone)
		assert.Equal(t, match.HAND, theirs.Zone)

		myHand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Contains(t, myHand, mine)
		assert.NotContains(t, myHand, theirs)
	})

	t.Run("the extra target may be left alone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		first := putCardInBattlezone(t, scn, opponent.Player, clonedSpiralCreatureUID, clonedSpiralSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, clonedSpiralCreatureUID, clonedSpiralSrc)

		_, err := player.Player.SpawnCard(clonedSpiralUID, match.GRAVEYARD)
		require.NoError(t, err)

		castSpell(t, scn, player, clonedSpiralUID)

		answerInTurn(t, scn, player, first.ID)
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, first.Zone)
		assert.Equal(t, match.BATTLEZONE, second.Zone)
	})

	t.Run("a lone creature is returned without asking", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		only := putCardInBattlezone(t, scn, opponent.Player, clonedSpiralCreatureUID, clonedSpiralSrc)

		castSpell(t, scn, player, clonedSpiralUID)
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, only.Zone)
	})

	t.Run("an empty board asks nothing", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		castSpell(t, scn, player, clonedSpiralUID)
		settleTurn(t, scn)
	})
}
