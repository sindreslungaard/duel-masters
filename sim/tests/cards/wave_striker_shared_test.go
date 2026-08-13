package cards

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/require"
)

// Shared plumbing for the wave striker creatures. Their abilities are only
// switched on while 2 or more *other* creatures in the battle zone have the
// keyword, so every behavioural test needs two spare wave strikers on the board.

const (
	// Macho Melon is the cheapest wave striker with no board effect of its own,
	// which makes it the ideal filler for making up the count.
	waveStrikerFillerUID = "fa987e39-2955-4074-bcf2-b7888ae27319"
	// A small blocker and a creature comfortably above the 2000 line, for the
	// cards that care about either.
	waveStrikerSmallBlockerUID = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (blocker, 2000)
	waveStrikerBigCreatureUID  = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	waveStrikerSharedSrc       = "wave_striker_shared_test_setup"
)

// spawnForLater puts a card in hand during setup so a following untap step
// gives it its keywords, the way a card drawn out of the deck would have them.
func spawnForLater(t *testing.T, player *match.PlayerReference, uid string) *match.Card {
	t.Helper()

	card, err := player.Player.SpawnCard(uid, match.HAND)
	require.NoError(t, err)

	return card
}

// putIntoPlay moves a card from hand into the battle zone and waits for the
// engine to settle.
func putIntoPlay(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, card *match.Card) *match.Card {
	t.Helper()

	moved, err := player.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, waveStrikerSharedSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	require.NoError(t, scn.WaitForEventLoop())

	return moved
}

// addWaveStrikerFillers puts n spare wave strikers into the player's battle
// zone. Two of them is what switches another wave striker's ability on.
func addWaveStrikerFillers(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, n int) {
	t.Helper()

	for range n {
		putCardInBattlezone(t, scn, player.Player, waveStrikerFillerUID, waveStrikerSharedSrc)
	}
}
