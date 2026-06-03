package cards

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const kingRippedHideUID = "f04feb7f-971f-4192-893a-46c23180233a"

// King Ripped-Hide lets the player draw up to 2 cards when summoned.
// Accepting the first prompt draws 1 card; the second prompt is then declined
// to avoid the ShowCardsNonDismissible preview after the second draw.
func TestKingRippedHide_AcceptsFirstDraw(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	// Keep the default deck intact so the player never decks out during the draw.
	// SpawnCard adds extra cards on top of the existing deck.
	const deckSeedUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	player.Player.SpawnCard(deckSeedUID, match.DECK)

	player.Player.SpawnCard(kingRippedHideUID, match.HAND)
	for range 7 {
		player.Player.SpawnCard(kingRippedHideUID, match.MANAZONE)
	}

	// Capture hand size after spawning so KRH is included in the baseline.
	handBefore, err := player.Player.Container(match.HAND)
	require.NoError(t, err)

	card, err := scn.FindCard(player.Player, match.HAND, kingRippedHideUID)
	require.NoError(t, err)

	// start captured before ActionPlayCard so Q1's prompt is inside the search window.
	start, err := scn.MessageCount(player)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	// Q1: wait for draw prompt, answer yes.
	require.NoError(t, scn.WaitForMessage(player, start, "action"))
	start, err = scn.MessageCount(player)
	require.NoError(t, err)
	require.NoError(t, scn.SubmitAction(player))

	// Dismiss the non-dismissible card preview that follows the first draw.
	require.NoError(t, scn.WaitForMessage(player, start, "show_cards_non_dismissible"))
	start, err = scn.MessageCount(player)
	require.NoError(t, err)
	require.NoError(t, scn.CancelAction(player)) // dismiss preview

	// Q2: wait for draw prompt then decline.
	require.NoError(t, scn.WaitForMessage(player, start, "action"))
	require.NoError(t, scn.CancelAction(player)) // decline second draw

	assert.Eventually(t, func() bool {
		handAfter, err := player.Player.Container(match.HAND)
		if err != nil {
			return false
		}
		// handBefore includes KRH. Played KRH (−1), drew 1 (+1) → net 0.
		return len(handAfter) == len(handBefore)
	}, time.Second, 10*time.Millisecond)
}

// Player declines the first draw prompt — the effect exits immediately with zero draws.
func TestKingRippedHide_DeclinesDraw(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.SpawnCard(kingRippedHideUID, match.HAND)
	for range 7 {
		player.Player.SpawnCard(kingRippedHideUID, match.MANAZONE)
	}

	// Capture hand size after spawning so KRH is included in the baseline.
	handBefore, err := player.Player.Container(match.HAND)
	require.NoError(t, err)

	card, err := scn.FindCard(player.Player, match.HAND, kingRippedHideUID)
	require.NoError(t, err)

	start, err := scn.MessageCount(player)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	// Wait for Q1 then decline — no second prompt follows.
	require.NoError(t, scn.WaitForMessage(player, start, "action"))
	require.NoError(t, scn.CancelAction(player))

	assert.Eventually(t, func() bool {
		handAfter, err := player.Player.Container(match.HAND)
		if err != nil {
			return false
		}
		// handBefore includes KRH. Played (−1), drew 0 → net −1.
		return len(handAfter) == len(handBefore)-1
	}, time.Second, 10*time.Millisecond)
}
