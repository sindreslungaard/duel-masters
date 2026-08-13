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
	gigavrandUID = "76310adb-f7c2-4545-8b71-2332b36fbb83"
	gigavrandSrc = "gigavrand_test_setup"
)

// gigavrandHandSnapshot copies a player's hand so a test can tell what happened
// to those exact cards. The turn that follows an end step deals a fresh card to
// whoever it belongs to, so counting the hand afterwards proves nothing.
func gigavrandHandSnapshot(t *testing.T, player *match.PlayerReference) []*match.Card {
	t.Helper()

	hand, err := player.Player.Container(match.HAND)
	require.NoError(t, err)
	require.NotEmpty(t, hand)

	snapshot := make([]*match.Card, len(hand))
	copy(snapshot, hand)

	return snapshot
}

// assertGigavrandHandKept checks that a snapshotted hand was left where it was.
func assertGigavrandHandKept(t *testing.T, snapshot []*match.Card) {
	t.Helper()

	for _, card := range snapshot {
		assert.Equal(t, match.HAND, card.Zone, card.Name+" should not have been discarded")
	}
}

// assertGigavrandHandDiscarded checks that a snapshotted hand was emptied.
func assertGigavrandHandDiscarded(t *testing.T, snapshot []*match.Card) {
	t.Helper()

	for _, card := range snapshot {
		assert.Equal(t, match.GRAVEYARD, card.Zone, card.Name+" should have been discarded")
	}
}

func TestGigavrand(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, gigavrandUID, gigavrandSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Gigavrand", 3000, 6, []string{civ.Darkness})
		assert.True(t, card.HasFamily(family.Chimera))
	})

	t.Run("two draws in a turn cost the opponent their hand", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, gigavrandUID, gigavrandSrc)

		// The turn's own draw is the first.
		require.NoError(t, scn.ActionEndTurn(player))

		// And this is the second.
		opponent.Player.DrawCards(1)

		require.NoError(t, scn.ActionEndTurn(opponent))
		settleTurn(t, scn)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Empty(t, hand)
	})

	t.Run("a single draw leaves the hand alone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, gigavrandUID, gigavrandSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		snapshot := gigavrandHandSnapshot(t, opponent)

		require.NoError(t, scn.ActionEndTurn(opponent))
		settleTurn(t, scn)

		assertGigavrandHandKept(t, snapshot)
	})

	t.Run("the count is per turn, not cumulative", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, gigavrandUID, gigavrandSrc)

		// One turn each way, with a single draw on each, is two draws in total
		// but never two in one turn.
		passTurnToSelf(t, scn, player, opponent)
		snapshot := gigavrandHandSnapshot(t, opponent)
		passTurnToSelf(t, scn, player, opponent)

		assertGigavrandHandKept(t, snapshot)
	})

	t.Run("it counts on its controller's turn too", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, gigavrandUID, gigavrandSrc)

		// The opponent draws nothing of their own during this turn, so both
		// draws have to be forced.
		opponent.Player.DrawCards(2)
		snapshot := gigavrandHandSnapshot(t, opponent)

		require.NoError(t, scn.ActionEndTurn(player))
		settleTurn(t, scn)

		assertGigavrandHandDiscarded(t, snapshot)
	})

	t.Run("its own controller drawing is not what it watches", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, gigavrandUID, gigavrandSrc)

		player.Player.DrawCards(3)
		snapshot := gigavrandHandSnapshot(t, opponent)

		require.NoError(t, scn.ActionEndTurn(player))
		settleTurn(t, scn)

		assertGigavrandHandKept(t, snapshot)
	})

	t.Run("out of the battle zone it punishes nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		gigavrand := putCardInBattlezone(t, scn, player.Player, gigavrandUID, gigavrandSrc)

		opponent.Player.DrawCards(2)
		snapshot := gigavrandHandSnapshot(t, opponent)

		_, err := player.Player.MoveCard(gigavrand.ID, match.BATTLEZONE, match.GRAVEYARD, gigavrandSrc)
		require.NoError(t, err)

		require.NoError(t, scn.ActionEndTurn(player))
		settleTurn(t, scn)

		assertGigavrandHandKept(t, snapshot)
	})

	t.Run("searching a card out of the deck is not drawing it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, gigavrandUID, gigavrandSrc)

		deck, err := opponent.Player.Container(match.DECK)
		require.NoError(t, err)

		for i := range 2 {
			_, err := opponent.Player.MoveCard(deck[i].ID, match.DECK, match.HAND, gigavrandSrc)
			require.NoError(t, err)
		}

		snapshot := gigavrandHandSnapshot(t, opponent)

		require.NoError(t, scn.ActionEndTurn(player))
		settleTurn(t, scn)

		assertGigavrandHandKept(t, snapshot)
	})
}
