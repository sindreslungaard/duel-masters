package cards

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Plumbing shared by every card test in this package. Anything tied to a single
// card or a single keyword belongs with that card instead.

const sharedTestSrc = "card_test_setup"

// setupDuel starts a match and returns the player whose turn it is plus their
// opponent.
func setupDuel(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}

// assertPrinted checks the numbers and civilizations on the card face. Keyword
// conditions are only rebuilt during an untap step, so callers that assert on
// those have to pass a turn first.
func assertPrinted(t *testing.T, card *match.Card, name string, power int, cost int, civs []string) {
	t.Helper()

	assert.Equal(t, name, card.Name)
	assert.Equal(t, power, card.Power)
	assert.Equal(t, cost, card.ManaCost)
	assert.Equal(t, civs, card.Civs)
	assert.Equal(t, civs, card.ManaRequirement)
}

// putCardInBattlezone spawns a card and moves it straight into play, skipping
// the summon. Note that a card conjured this way has no conditions until an
// untap step runs, so a test that depends on its keywords has to pass a turn.
func putCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string, source string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, source)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}

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
//
// A card whose arrival opens a prompt cannot use this: the move runs on the
// test goroutine, which would then be the only one able to answer a prompt it
// is itself blocked on. Summon those through scn.ActionPlayCard instead.
func putIntoPlay(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, card *match.Card) *match.Card {
	t.Helper()

	moved, err := player.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, sharedTestSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	require.NoError(t, scn.WaitForEventLoop())

	return moved
}

// castSpell plays a spell from hand, paying for it with copies of itself so the
// civilization always matches.
func castSpell(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, uid string) *match.Card {
	t.Helper()

	spell, err := player.Player.SpawnCard(uid, match.HAND)
	require.NoError(t, err)

	card, err := scn.FindCard(player.Player, match.HAND, uid)
	require.NoError(t, err)

	for range card.ManaCost {
		_, err := player.Player.SpawnCard(uid, match.MANAZONE)
		require.NoError(t, err)
	}

	require.NoError(t, scn.ActionPlayCard(player, spell.ID))

	return spell
}

// passTurnToSelf ends both players' turns, leaving the match at the start of a
// fresh turn for player. A start-of-turn prompt, if any, is open when it
// returns, because the turn transition blocks on it.
func passTurnToSelf(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, opponent *match.PlayerReference) {
	t.Helper()

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))
}

// settleTurn waits for the rest of the turn to finish after the last prompt of
// a sequence has been answered.
func settleTurn(t *testing.T, scn *scenario.TestScenario) {
	t.Helper()

	require.NoError(t, scn.WaitForEventLoop())
}

// answerInTurn answers one prompt and waits for the engine's reply, so a
// follow-up answer is never sent while the previous one is still in flight.
//
// This matters because the parallel action handler drains Player.Action before
// writing to it, to stop a client flooding the match with goroutines. An answer
// that has not been picked up yet is thrown away by the next answer's drain, so
// two answers sent back to back can lose the first one and leave the event loop
// waiting on a prompt nobody will answer again. A real client cannot do this,
// because it only ever shows one prompt at a time.
func answerInTurn(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, cardIDs ...string) {
	t.Helper()

	start, err := scn.MessageCount(player)
	require.NoError(t, err)

	require.NoError(t, scn.SubmitAction(player, cardIDs...))
	require.NoError(t, scn.WaitForMessage(player, start, "action", "action_error", "warn", "wait", "state_update"))
}

// cancelInTurn is answerInTurn for a cancellation.
func cancelInTurn(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) {
	t.Helper()

	start, err := scn.MessageCount(player)
	require.NoError(t, err)

	require.NoError(t, scn.CancelAction(player))
	require.NoError(t, scn.WaitForMessage(player, start, "action", "action_error", "warn", "wait", "state_update"))
}

// answerDrawUpTo answers the prompts fx.DrawUpto opens. It asks once per card
// rather than once in total, and each accepted draw opens a preview of the card
// that has to be acknowledged before the effect continues.
func answerDrawUpTo(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, max int, take bool) {
	t.Helper()

	if !take {
		// Declining the first question ends the whole draw.
		cancelInTurn(t, scn, player)
		return
	}

	for range max {
		// Yes to this card, then close the preview it opens. The preview is
		// dismissed with a cancel, which is what its Close button sends.
		answerInTurn(t, scn, player)
		cancelInTurn(t, scn, player)
	}
}

// answerOrderPrompt answers an fx.OrderCards prompt by keeping the order the
// engine offered. The prompt only opens for two or more cards.
func answerOrderPrompt(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) {
	t.Helper()

	action, err := scn.LatestAction(player, 0)
	require.NoError(t, err, "expected an order prompt to be open")
	require.GreaterOrEqual(t, len(action.Cards), 2)

	ids := make([]string, 0, len(action.Cards))
	for _, card := range action.Cards {
		ids = append(ids, card.CardID)
	}

	answerInTurn(t, scn, player, ids...)
}

// emptyHand clears a player's hand so a test can seed exactly what it needs.
func emptyHand(t *testing.T, player *match.PlayerReference, source string) {
	t.Helper()

	hand, err := player.Player.Container(match.HAND)
	require.NoError(t, err)

	for _, card := range hand {
		_, err := player.Player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, source)
		require.NoError(t, err)
	}
}

// countHeaders counts how many messages of one kind a player received.
func countHeaders(headers []string, wanted string) int {
	count := 0
	for _, header := range headers {
		if header == wanted {
			count++
		}
	}

	return count
}

// summonWithOwnMana puts a card in hand with enough copies of itself in the
// mana zone to pay for it, then summons it through the event loop. A card whose
// arrival opens a prompt has to be summoned this way rather than moved with
// putIntoPlay, which would leave the test goroutine blocked on a prompt only it
// could answer.
func summonWithOwnMana(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, uid string) *match.Card {
	t.Helper()

	card, err := player.Player.SpawnCard(uid, match.HAND)
	require.NoError(t, err)

	for range card.ManaCost {
		_, err := player.Player.SpawnCard(uid, match.MANAZONE)
		require.NoError(t, err)
	}

	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	return card
}

// discardTriggerUID is a creature whose summon makes the opponent discard a
// card at random, which is how a test triggers a discard replacement.
const discardTriggerUID = "ba955ab0-5bb3-4aaf-82f3-293522e65a9c" // Locomotiver

// discardTriggerManaUID pays for it.
const discardTriggerManaUID = "e2b992ee-91a3-49d3-8228-7be60a0b9ec5" // Writhing Bone Ghoul (darkness)

// setupDiscardReplacementTest leaves the opponent on turn with a creature ready
// to summon that discards a random card, and the card under test as the only
// card in the other player's hand so the discard can only hit that.
func setupDiscardReplacementTest(t *testing.T, uid string) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference, *match.Card) {
	t.Helper()

	scn, player, opponent := setupDuel(t)

	for range 4 {
		_, err := opponent.Player.SpawnCard(discardTriggerManaUID, match.MANAZONE)
		require.NoError(t, err)
	}

	_, err := opponent.Player.SpawnCard(discardTriggerUID, match.HAND)
	require.NoError(t, err)

	require.NoError(t, scn.ActionEndTurn(player))
	require.True(t, scn.Match.IsPlayerTurn(opponent.Player))

	// Emptied last so the drawn opening hand cannot be discarded instead.
	emptyHand(t, player, sharedTestSrc)

	card, err := player.Player.SpawnCard(uid, match.HAND)
	require.NoError(t, err)

	return scn, player, opponent, card
}

// discardTriggerCard finds the creature setupDiscardReplacementTest prepared.
func discardTriggerCard(t *testing.T, scn *scenario.TestScenario, opponent *match.PlayerReference) string {
	t.Helper()

	card, err := scn.FindCard(opponent.Player, match.HAND, discardTriggerUID)
	require.NoError(t, err)

	return card.ID
}
