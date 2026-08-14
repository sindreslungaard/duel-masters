package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	wilyCarpenterUID = "6fefa2b0-3d88-4420-88b3-9896b6ac4ce5"
)

func TestWilyCarpenter(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		// Kept in hand: it asks how much to draw as it arrives, and a card moved into play from the test
		// goroutine would leave it waiting on a prompt only it could answer.
		card := spawnForLater(t, player, wilyCarpenterUID)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Wily Carpenter", 1000, 3, []string{civ.Water})
		assert.True(t, card.HasFamily(family.Merfolk))
	})

	t.Run("draw two then discard two is a net loss of nothing", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, wilyCarpenterUID)

		answerDrawUpTo(t, scn, player, 2, true)
		answerInTurn(t, scn, player, handBefore[0].ID, handBefore[1].ID)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		// The carpenter left the hand to be summoned, two were drawn and two
		// discarded.
		assert.Len(t, hand, len(handBefore))
		assert.Equal(t, match.GRAVEYARD, handBefore[0].Zone)
		assert.Equal(t, match.GRAVEYARD, handBefore[1].Zone)
	})
}
