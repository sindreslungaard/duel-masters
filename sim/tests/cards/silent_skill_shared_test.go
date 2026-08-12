package cards

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

// Shared plumbing for the silent skill creatures. The keyword resolves at the
// start of its controller's turn, so every one of these tests has to hand the
// turn over and take it back before the ability is offered at all.

func setupSilentSkillTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}

// passTurnToSelf ends both players' turns, leaving the match at the start of a
// fresh turn for player. A silent skill prompt, if any, is open when it
// returns, because the turn transition blocks on it.
func passTurnToSelf(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, opponent *match.PlayerReference) {
	t.Helper()

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))
}

// useSilentSkill accepts the pending silent skill prompt. Any follow-up prompt
// the ability itself opens is left for the caller to answer.
func useSilentSkill(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) {
	t.Helper()

	answerSilentSkillPrompt(t, scn, player, true)
}

// declineSilentSkill answers no, which untaps the creature as a normal untap
// step would have.
func declineSilentSkill(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference) {
	t.Helper()

	answerSilentSkillPrompt(t, scn, player, false)
}

func answerSilentSkillPrompt(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, use bool) {
	t.Helper()

	start, err := scn.MessageCount(player)
	require.NoError(t, err)

	action, err := scn.LatestAction(player, 0)
	require.NoError(t, err, "expected a silent skill prompt to be open")
	require.NotEmpty(t, action.Text)

	if use {
		require.NoError(t, scn.SubmitAction(player))
	} else {
		require.NoError(t, scn.CancelAction(player))
	}

	require.NoError(t, scn.WaitForMessage(player, start, "action", "wait", "state_update"))

	// The ability may have opened a prompt of its own, and the event loop stays
	// blocked until that one is answered, so there is nothing to settle yet. A
	// "wait" means the prompt went to the opponent instead, which blocks the
	// loop just the same.
	headers, err := scn.MessageHeaders(player, start)
	require.NoError(t, err)
	if slices.Contains(headers, "action") || slices.Contains(headers, "wait") {
		return
	}

	require.NoError(t, scn.WaitForEventLoop())
}

// settleTurn waits for the rest of the turn to finish after the last prompt of
// a silent skill sequence has been answered.
func settleTurn(t *testing.T, scn *scenario.TestScenario) {
	t.Helper()

	require.NoError(t, scn.WaitForEventLoop())
}
