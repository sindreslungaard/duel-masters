package tests

import (
	"duel-masters/tests/scenario"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const spectatorTestTimeout = 5 * time.Second

// Sending the match state to a spectator whose connection is slow must not stop
// other spectators from joining, so the spectator list cannot stay locked while
// state is being sent.
func TestSpectatorCanJoinWhileStateIsSentToSlowSpectator(t *testing.T) {
	scn := scenario.New()

	_, slowConn, err := scn.JoinSpectator("spectator-slow", "SlowSpectator")
	if err != nil {
		t.Fatal(err)
	}
	if err := scn.WaitForSpectatorMessage(slowConn, 0, "state_update"); err != nil {
		t.Fatal(err)
	}

	writeHeld := slowConn.HoldWrites()
	t.Cleanup(slowConn.ReleaseWrites)

	broadcastDone := make(chan struct{})
	go func() {
		defer close(broadcastDone)
		scn.Match.BroadcastState()
	}()

	select {
	case <-writeHeld:
	case <-time.After(spectatorTestTimeout):
		t.Fatal("the state was never sent to the slow spectator")
	}

	since, err := scn.MessageCount(scn.Match.Player1)
	if err != nil {
		t.Fatal(err)
	}

	_, joiningConn, err := scn.JoinSpectator("spectator-joining", "JoiningSpectator")
	if err != nil {
		t.Fatal(err)
	}

	// A join announces itself in the chat after the spectator is registered, and
	// the players receive the chat before any spectator does.
	if err := scn.WaitForMessage(scn.Match.Player1, since, "chat"); err != nil {
		t.Fatalf("a spectator could not join while the state was being sent to a slow spectator: %v", err)
	}

	chat, err := scn.ChatMessages(scn.Match.Player1, since)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, chat, "JoiningSpectator started spectating")

	slowConn.ReleaseWrites()

	select {
	case <-broadcastDone:
	case <-time.After(spectatorTestTimeout):
		t.Fatal("broadcasting the state did not finish after the slow spectator caught up")
	}

	if err := scn.WaitForSpectatorMessage(joiningConn, 0, "state_update"); err != nil {
		t.Fatalf("the joining spectator never received the match state: %v", err)
	}
}
