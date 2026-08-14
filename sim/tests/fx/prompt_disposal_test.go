package fx

import (
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// immortalBaronVorgUID is a 2 mana fire creature, used here only because playing
// it opens the mana payment prompt.
const immortalBaronVorgUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5"

// promptReturnTimeout is generous on purpose. The prompts under test either
// abort promptly once the action channel is closed, or they never return at
// all, so the exact value only decides how long a regression takes to report.
const promptReturnTimeout = 2 * time.Second

// TestPromptsAbandonThemselvesOnDisposal covers the prompts that block the match
// event loop while waiting for a player.
//
// Disposing a match closes both players' action channels. A receive from a
// closed channel returns the zero PlayerAction immediately and forever, so a
// prompt that treats it as an invalid answer and retries turns the match event
// loop into a busy loop that outlives the match and burns a core until the
// process restarts. Every prompt must instead abort the effect in progress, so
// that no card logic resolves for a game that no longer exists.
func TestPromptsAbandonThemselvesOnDisposal(t *testing.T) {

	prompts := []struct {
		name string
		// header is the message the prompt sends to the player, waited for so
		// the match is only disposed once the prompt is actually open.
		header string
		prompt func(*scenario.TestScenario, *match.Player)
	}{
		{
			name:   "card selection",
			header: "action",
			prompt: func(scn *scenario.TestScenario, p *match.Player) {
				fx.Select(p, scn.Match, p, match.HAND, "Select a card", 1, 1, false)
			},
		},
		{
			name:   "count selection",
			header: "action",
			prompt: func(scn *scenario.TestScenario, p *match.Player) {
				fx.SelectCount(p, scn.Match, "How many?", 1, 3)
			},
		},
		{
			name:   "yes/no question",
			header: "action",
			prompt: func(scn *scenario.TestScenario, p *match.Player) {
				fx.BinaryQuestion(p, scn.Match, "Do you want to?")
			},
		},
		{
			name:   "card ordering",
			header: "action",
			prompt: func(scn *scenario.TestScenario, p *match.Player) {
				hand, err := p.Container(match.HAND)
				require.NoError(t, err)
				fx.OrderCards(p, scn.Match, hand, "Order these")
			},
		},
		{
			name:   "multiple choice",
			header: "action",
			prompt: func(scn *scenario.TestScenario, p *match.Player) {
				fx.MultipleChoiceQuestion(p, scn.Match, "Pick one", []string{"a", "b", "c"})
			},
		},
		{
			name:   "shown cards that must be dismissed",
			header: "show_cards_non_dismissible",
			prompt: func(scn *scenario.TestScenario, p *match.Player) {
				scn.Match.ShowCardsNonDismissible(p, "Look at these", []string{"some-image-id"})
			},
		},
	}

	for _, prompt := range prompts {
		t.Run(prompt.name, func(t *testing.T) {
			scn := scenario.New()
			player := scn.Match.Player1

			start, err := scn.MessageCount(player)
			require.NoError(t, err)

			aborted := make(chan any, 1)
			go func() {
				// The prompt aborts by panicking, which is how it unwinds every
				// enclosing card handler. Nothing else in this goroutine may
				// recover it, or the abort would be reported as a fault.
				defer func() { aborted <- recover() }()
				prompt.prompt(scn, player.Player)
				aborted <- nil
			}()

			require.NoError(t, scn.WaitForMessage(player, start, prompt.header), "the prompt should have been opened")

			scn.Match.Dispose()

			select {
			case recovered := <-aborted:
				assert.True(t, match.IsMatchDisposed(recovered),
					"the prompt must abort with the match disposed signal rather than resolving, got %v", recovered)
			case <-time.After(promptReturnTimeout):
				t.Fatal("the prompt never returned after the match was disposed; it is spinning on the closed action channel")
			}
		})
	}
}

// TestDisconnectDuringPromptStopsMatch drives the production path end to end:
// a card is played through Match.Parse so the mana payment prompt blocks the
// match event loop, then both websockets close, which is what makes
// OnSocketClose dispose the match and close the players' action channels.
//
// It also covers the abort reaching its recover: if the match event loop did not
// recover the abort raised by the abandoned prompt, it would take the test
// process down with it.
func TestDisconnectDuringPromptStopsMatch(t *testing.T) {
	scn := scenario.New()

	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.Player1
	if opponent == player {
		opponent = scn.Match.Player2
	}

	for range 3 {
		_, err := player.Player.SpawnCard(immortalBaronVorgUID, match.MANAZONE)
		require.NoError(t, err)
	}

	creature, err := player.Player.SpawnCard(immortalBaronVorgUID, match.HAND)
	require.NoError(t, err)

	start, err := scn.MessageCount(player)
	require.NoError(t, err)

	scn.Match.Parse(player.Socket, []byte(fmt.Sprintf(`{"header":"add_to_playzone","virtualId":%q}`, creature.ID)))

	// The match event loop is now blocked inside the mana payment prompt
	require.NoError(t, scn.WaitForMessage(player, start, "action"))

	// Both players disconnect while that prompt is still open, which disposes
	// the match from the socket goroutine
	player.Socket.Close()
	opponent.Socket.Close()

	require.Eventually(t, scn.Match.IsClosed, promptReturnTimeout, 10*time.Millisecond,
		"the second disconnect should dispose the match")

	require.Eventually(t, scn.Match.EventLoopStopped, promptReturnTimeout, 10*time.Millisecond,
		"the match event loop must unwind out of the abandoned prompt and stop")

	assert.Equal(t, match.HAND, creature.Zone, "an abandoned payment must not summon the creature")
}
