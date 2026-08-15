package match

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDeckStringUsesCreatedCardIDs(t *testing.T) {
	player := NewPlayer(&Match{}, 1)
	player.deck = []*Card{{ImageID: "card-a"}, {ImageID: "card-b"}, {ImageID: "card-a"}}

	if deck := player.deckString(); deck != "card-a,card-b,card-a" {
		t.Fatalf("expected created card IDs in deck string, got %q", deck)
	}
}

func TestNewDuelResultWebhookPayload(t *testing.T) {
	m := &Match{
		ID:          "duel-123",
		MatchName:   "Ranked Duel",
		FormatID:    "standard:11111111-1111-4111-8111-111111111111",
		FormatName:  "Classic",
		turnsPlayed: 14,
		Player1:     &PlayerReference{UID: "host-1", Username: "Alice", DeckStr: "card-a,card-b"},
		Player2:     &PlayerReference{UID: "guest-2", Username: "Bob", DeckStr: "card-c,card-d"},
	}

	duel := DuelRecord{
		UID:             "duel-123",
		Host:            "host-1",
		HostDeck:        "card-a,card-b",
		Guest:           "guest-2",
		GuestDeck:       "card-c,card-d",
		Started:         100,
		Ended:           190,
		Turns:           14,
		Winner:          "guest-2",
		WonByDisconnect: true,
	}

	payload := m.newDuelResultWebhookPayload(duel)

	if payload.MatchID != "duel-123" {
		t.Fatalf("expected match id to be preserved, got %q", payload.MatchID)
	}

	if payload.FormatID != "standard:11111111-1111-4111-8111-111111111111" ||
		payload.FormatName != "Classic" {
		t.Fatalf("expected the format descriptor to be reported, got %+v", payload)
	}

	if payload.DurationSeconds != 90 {
		t.Fatalf("expected duration to be 90 seconds, got %d", payload.DurationSeconds)
	}

	if payload.Turns != 14 {
		t.Fatalf("expected turns to be 14, got %d", payload.Turns)
	}

	if payload.Winner == nil || payload.Winner.UserID != "guest-2" {
		t.Fatalf("expected winner to be guest-2, got %+v", payload.Winner)
	}

	if payload.Loser == nil || payload.Loser.UserID != "host-1" {
		t.Fatalf("expected loser to be host-1, got %+v", payload.Loser)
	}

	if payload.Host == nil || payload.Host.Username != "Alice" {
		t.Fatalf("expected host info to be included, got %+v", payload.Host)
	}

	if payload.Guest == nil || payload.Guest.Username != "Bob" {
		t.Fatalf("expected guest info to be included, got %+v", payload.Guest)
	}

	if payload.Host.Deck != "card-a,card-b" {
		t.Fatalf("expected host's actual deck cards, got %#v", payload.Host.Deck)
	}

	if payload.Guest.Deck != "card-c,card-d" {
		t.Fatalf("expected guest's actual deck cards, got %#v", payload.Guest.Deck)
	}

	if payload.StartedAt != time.Unix(100, 0).UTC().Format(time.RFC3339) {
		t.Fatalf("expected RFC3339 start time, got %q", payload.StartedAt)
	}

	if payload.EndedAt != time.Unix(190, 0).UTC().Format(time.RFC3339) {
		t.Fatalf("expected RFC3339 end time, got %q", payload.EndedAt)
	}
}

func TestSendDuelResultWebhookMatchesShobuContract(t *testing.T) {
	var gotAuthorization string
	var gotContentType string
	var got map[string]any
	var decodeErr error

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuthorization = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		decodeErr = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	payload := duelResultWebhookPayload{
		MatchID: "duel-123",
		EndedAt: "2026-08-08T01:02:03Z",
		Host:    &duelResultWebhookPlayer{UserID: "host-1", Username: "Alice", Deck: "card-a,card-b"},
		Guest:   &duelResultWebhookPlayer{UserID: "guest-2", Username: "Bob", Deck: "card-c,card-d"},
	}

	if err := sendDuelResultWebhook(server.URL, " duel-secret ", payload); err != nil {
		t.Fatalf("expected webhook request to succeed, got %v", err)
	}

	if gotAuthorization != "Bearer duel-secret" {
		t.Fatalf("expected bearer authorization header, got %q", gotAuthorization)
	}

	if gotContentType != "application/json" {
		t.Fatalf("expected JSON content type, got %q", gotContentType)
	}

	if decodeErr != nil {
		t.Fatalf("failed to decode request body: %v", decodeErr)
	}

	if got["matchId"] != "duel-123" {
		t.Fatalf("expected Shobu matchId field, got %#v", got["matchId"])
	}

	host, ok := got["host"].(map[string]any)
	if !ok || host["userId"] != "host-1" || host["username"] != "Alice" {
		t.Fatalf("expected Shobu host participant fields, got %#v", got["host"])
	}

	if host["deck"] != "card-a,card-b" {
		t.Fatalf("expected host deck string, got %#v", host["deck"])
	}

	if got["endedAt"] != "2026-08-08T01:02:03Z" {
		t.Fatalf("expected Shobu endedAt field, got %#v", got["endedAt"])
	}

	for _, legacyField := range []string{"duel_id", "ended_at"} {
		if _, exists := got[legacyField]; exists {
			t.Fatalf("legacy field %q must not be sent", legacyField)
		}
	}
}
