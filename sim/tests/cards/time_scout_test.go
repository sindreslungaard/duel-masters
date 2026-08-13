package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	timeScoutUID      = "23c905d3-a554-4d45-b8d7-1b46fe8d117f"
	timeScoutSetupSrc = "time_scout_test_setup"
)

func TestTimeScout(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		scout := putCardInBattlezone(t, scn, player.Player, timeScoutUID, timeScoutSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, scout, "Time Scout", 1000, 2, []string{civ.Water})
		assert.True(t, scout.HasFamily(family.Merfolk))
	})

	t.Run("it shows its controller the top card and leaves it there", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		scout, err := player.Player.SpawnCard(timeScoutUID, match.HAND)
		require.NoError(t, err)
		for range 2 {
			_, err := player.Player.SpawnCard(timeScoutUID, match.MANAZONE)
			require.NoError(t, err)
		}

		deckBefore, err := opponent.Player.Container(match.DECK)
		require.NoError(t, err)
		top := opponent.Player.PeekDeck(1)
		require.Len(t, top, 1)

		playerStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionPlayCard(player, scout.ID))
		require.NoError(t, scn.WaitForEventLoop())

		playerHeaders, err := scn.MessageHeaders(player, playerStart)
		require.NoError(t, err)
		assert.Contains(t, playerHeaders, "show_cards")

		opponentHeaders, err := scn.MessageHeaders(opponent, opponentStart)
		require.NoError(t, err)
		assert.NotContains(t, opponentHeaders, "show_cards", "only the controller looks")

		deck, err := opponent.Player.Container(match.DECK)
		require.NoError(t, err)
		assert.Len(t, deck, len(deckBefore), "looking does not move the card")
		assert.Equal(t, match.DECK, top[0].Zone)
	})
}
