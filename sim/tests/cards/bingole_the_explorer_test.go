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
	bingoleTheExplorerUID = "ba6ad07e-eed3-49b9-8830-2df794b13066"
	bingoleSetupSrc       = "bingole_the_explorer_test_setup"
)

func TestBingoleTheExplorer(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, bingoleTheExplorerUID, bingoleSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Bingole, the Explorer", 3000, 4, []string{civ.Light})
		assert.True(t, card.HasFamily(family.Gladiator))
	})

	t.Run("a discard on the opponent's turn may put it into play instead", func(t *testing.T) {
		scn, player, opponent, bingole := setupDiscardReplacementTest(t, bingoleTheExplorerUID)

		require.NoError(t, scn.ActionPlayCard(opponent, discardTriggerCard(t, scn, opponent)))
		require.NoError(t, scn.SubmitAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, bingole.Zone)
	})

	t.Run("declining discards it as normal", func(t *testing.T) {
		scn, player, opponent, bingole := setupDiscardReplacementTest(t, bingoleTheExplorerUID)

		require.NoError(t, scn.ActionPlayCard(opponent, discardTriggerCard(t, scn, opponent)))
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, bingole.Zone)
	})
}
