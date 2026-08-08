package match

import (
	"testing"
)

func TestNewDuelFinishedMessageIncludesWinnerAndGeneratedFlag(t *testing.T) {
	m := &Match{
		ID:          "duel-123",
		turnsPlayed: 14,
		Player1:     &PlayerReference{UID: "host-1", Username: "Alice"},
		Player2:     &PlayerReference{UID: "guest-2", Username: "Bob"},
	}

	duel := DuelRecord{
		UID:             "duel-123",
		Started:         100,
		Ended:           190,
		Turns:           14,
		Winner:          "guest-2",
		WonByDisconnect: true,
	}

	message := m.newDuelFinishedMessage(duel, true)

	if message.Header != "duel_finished" {
		t.Fatalf("expected duel_finished header, got %q", message.Header)
	}

	if message.Winner == nil || message.Winner.UID != "guest-2" || message.Winner.Username != "Bob" {
		t.Fatalf("expected winner to be Bob/guest-2, got %+v", message.Winner)
	}

	if !message.MatchResultGenerated {
		t.Fatal("expected match result generated to be true")
	}

	if !message.WonByDisconnect {
		t.Fatal("expected won by disconnect to be true")
	}

	if message.DurationSeconds != 90 {
		t.Fatalf("expected duration to be 90, got %d", message.DurationSeconds)
	}

	if message.Turns != 14 {
		t.Fatalf("expected turns to be 14, got %d", message.Turns)
	}
}

func TestShouldGenerateMatchResultDefaultsToOneMinute(t *testing.T) {
	t.Setenv("duel_result_min_duration_seconds", "")

	if shouldGenerateMatchResult(DuelRecord{Started: 100, Ended: 159}) {
		t.Fatal("expected short duel to be ignored")
	}

	if !shouldGenerateMatchResult(DuelRecord{Started: 100, Ended: 160}) {
		t.Fatal("expected 60-second duel to count")
	}
}

func TestShouldGenerateMatchResultUsesConfiguredMinimumDuration(t *testing.T) {
	t.Setenv("duel_result_min_duration_seconds", "10")

	if shouldGenerateMatchResult(DuelRecord{Started: 100, Ended: 109}) {
		t.Fatal("expected duel below configured duration to be ignored")
	}

	if !shouldGenerateMatchResult(DuelRecord{Started: 100, Ended: 110}) {
		t.Fatal("expected duel at configured duration to count")
	}
}

func TestShouldGenerateMatchResultFallsBackToOneMinuteForInvalidConfiguration(t *testing.T) {
	t.Setenv("duel_result_min_duration_seconds", "invalid")

	if shouldGenerateMatchResult(DuelRecord{Started: 100, Ended: 159}) {
		t.Fatal("expected invalid configuration to use the 60-second default")
	}

	if !shouldGenerateMatchResult(DuelRecord{Started: 100, Ended: 160}) {
		t.Fatal("expected invalid configuration to use the 60-second default")
	}
}
