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
	funkyWizardUID = "3b11f35d-0e89-49f9-bcc8-0fa91a19097d"
)

func TestFunkyWizard(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		// Kept in hand: both players are asked whether to draw as it arrives, and a card moved into play from the test
		// goroutine would leave it waiting on a prompt only it could answer.
		card := spawnForLater(t, player, funkyWizardUID)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Funky Wizard", 2000, 4, []string{civ.Water})
		assert.True(t, card.HasFamily(family.Merfolk))
	})

	t.Run("both players may draw", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		myHandBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		theirHandBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, funkyWizardUID)

		answerInTurn(t, scn, player)
		answerInTurn(t, scn, opponent)
		require.NoError(t, scn.WaitForEventLoop())

		myHand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		theirHand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		// The wizard was spawned into hand and then left it to be summoned, so
		// the only lasting change is the card it drew.
		assert.Len(t, myHand, len(myHandBefore)+1)
		assert.Len(t, theirHand, len(theirHandBefore)+1)
	})

	t.Run("either player may decline", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		theirHandBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, funkyWizardUID)

		cancelInTurn(t, scn, player)
		cancelInTurn(t, scn, opponent)
		require.NoError(t, scn.WaitForEventLoop())

		theirHand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, theirHand, len(theirHandBefore))
	})
}
