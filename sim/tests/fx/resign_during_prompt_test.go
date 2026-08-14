package fx

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResignWhileAPromptIsPending covers the case a player is most likely to
// want to leave in: the match event loop is blocked waiting for somebody to
// answer a prompt.
//
// Resigning has to reach the match without going through the loop, because the
// loop cannot process anything while it waits. A sequential resign would sit in
// the queue forever and neither player could get out of the duel.
func TestResignWhileAPromptIsPending(t *testing.T) {
	for _, resigner := range []string{"the player being prompted", "the opponent"} {
		t.Run(resigner, func(t *testing.T) {
			scn := scenarioWithPendingPrompt(t)

			prompted := scn.Match.CurrentPlayer()
			opponent := scn.Match.Player1
			if opponent == prompted {
				opponent = scn.Match.Player2
			}

			socket := prompted.Socket
			if resigner == "the opponent" {
				socket = opponent.Socket
			}

			scn.Match.Parse(socket, []byte(`{"header":"resign"}`))

			require.Eventually(t, scn.Match.IsClosed, promptReturnTimeout, 10*time.Millisecond,
				"resigning must end the match even though the event loop is waiting on a prompt")

			require.Eventually(t, scn.Match.EventLoopStopped, promptReturnTimeout, 10*time.Millisecond,
				"the event loop must unwind out of the abandoned prompt and stop")
		})
	}
}

// TestResignEndsTheMatchWithNoPromptPending is the ordinary case, kept so that
// moving the resign off the event loop cannot quietly break it.
func TestResignEndsTheMatchWithNoPromptPending(t *testing.T) {
	scn := scenario.New()

	scn.Match.Parse(scn.Match.Player1.Socket, []byte(`{"header":"resign"}`))

	require.Eventually(t, scn.Match.IsClosed, promptReturnTimeout, 10*time.Millisecond,
		"resigning should end the match")
}

// scenarioWithPendingPrompt leaves the match event loop blocked inside the mana
// payment prompt of a card that is being played.
func scenarioWithPendingPrompt(t *testing.T) *scenario.TestScenario {
	t.Helper()

	scn := scenario.New()

	player := scn.Match.CurrentPlayer()

	for range 3 {
		_, err := player.Player.SpawnCard(immortalBaronVorgUID, match.MANAZONE)
		require.NoError(t, err)
	}

	creature, err := player.Player.SpawnCard(immortalBaronVorgUID, match.HAND)
	require.NoError(t, err)

	start, err := scn.MessageCount(player)
	require.NoError(t, err)

	scn.Match.Parse(player.Socket, []byte(fmt.Sprintf(`{"header":"add_to_playzone","virtualId":%q}`, creature.ID)))

	require.NoError(t, scn.WaitForMessage(player, start, "action"),
		"the mana payment prompt should be open and holding the event loop")

	assert.Equal(t, match.HAND, creature.Zone)

	return scn
}
