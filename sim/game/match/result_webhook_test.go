package match

import (
	"duel-masters/db"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewDuelResultWebhookPayload(t *testing.T) {
	m := &Match{
		ID:          "duel-123",
		MatchName:   "Ranked Duel",
		Format:      RegularFormat,
		turnsPlayed: 14,
		Player1:     &PlayerReference{UID: "host-1", Username: "Alice", DeckStr: "deck-a"},
		Player2:     &PlayerReference{UID: "guest-2", Username: "Bob", DeckStr: "deck-b"},
	}

	duel := db.Duel{
		UID:             "duel-123",
		Format:          string(RegularFormat),
		Host:            "host-1",
		HostDeck:        "deck-a",
		Guest:           "guest-2",
		GuestDeck:       "deck-b",
		Started:         100,
		Ended:           190,
		Turns:           14,
		Winner:          "guest-2",
		WonByDisconnect: true,
	}

	payload := m.newDuelResultWebhookPayload(duel)

	if payload.DuelID != "duel-123" {
		t.Fatalf("expected duel id to be preserved, got %q", payload.DuelID)
	}

	if payload.DurationSeconds != 90 {
		t.Fatalf("expected duration to be 90 seconds, got %d", payload.DurationSeconds)
	}

	if payload.Turns != 14 {
		t.Fatalf("expected turns to be 14, got %d", payload.Turns)
	}

	if payload.Winner == nil || payload.Winner.UID != "guest-2" {
		t.Fatalf("expected winner to be guest-2, got %+v", payload.Winner)
	}

	if payload.Loser == nil || payload.Loser.UID != "host-1" {
		t.Fatalf("expected loser to be host-1, got %+v", payload.Loser)
	}

	if payload.Host == nil || payload.Host.Username != "Alice" {
		t.Fatalf("expected host info to be included, got %+v", payload.Host)
	}

	if payload.Guest == nil || payload.Guest.Username != "Bob" {
		t.Fatalf("expected guest info to be included, got %+v", payload.Guest)
	}
}

func TestSendDuelResultWebhookSetsAuthorizationHeader(t *testing.T) {
	var gotAuthorization string
	var gotPayload duelResultWebhookPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuthorization = r.Header.Get("Authorization")

		if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
			t.Fatalf("failed to decode request payload: %v", err)
		}

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	payload := duelResultWebhookPayload{DuelID: "duel-123", Turns: 7}

	if err := sendDuelResultWebhook(server.URL, "duel-secret", payload); err != nil {
		t.Fatalf("expected webhook request to succeed, got %v", err)
	}

	if gotAuthorization != "duel-secret" {
		t.Fatalf("expected Authorization header to contain duel secret, got %q", gotAuthorization)
	}

	if gotPayload.DuelID != "duel-123" || gotPayload.Turns != 7 {
		t.Fatalf("unexpected webhook payload received: %+v", gotPayload)
	}
}
