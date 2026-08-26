package api

import (
	"bytes"
	"duel-masters/game/match"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetMatchesHandlerRequiresServerToken(t *testing.T) {
	t.Setenv("secret", "test-secret")

	api := New(match.NewSystem())

	for _, test := range []struct {
		name          string
		authorization string
	}{
		{name: "missing token"},
		{name: "wrong token", authorization: "Bearer wrong-secret"},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/matches", nil)
			req.Header.Set("Authorization", test.authorization)
			res := httptest.NewRecorder()

			api.getMatchesHandler(res, req)

			if res.Code != http.StatusUnauthorized {
				t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
			}
		})
	}
}

func TestGetMatchesHandlerRequiresConfiguredSecret(t *testing.T) {
	t.Setenv("secret", "")

	api := New(match.NewSystem())
	req := httptest.NewRequest(http.MethodGet, "/api/matches", nil)
	req.Header.Set("Authorization", "Bearer any-secret")
	res := httptest.NewRecorder()

	api.getMatchesHandler(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestGetMatchesHandlerReturnsAllCurrentMatches(t *testing.T) {
	t.Setenv("secret", "test-secret")

	system := match.NewSystem()
	m := system.NewMatch(
		"Private duel",
		"host-1",
		"Alice",
		[]string{"secret-host-card"},
		"guest-1",
		"Bob",
		[]string{"secret-guest-card"},
		false,
		false,
		match.FormatDescriptor{
			ID:     "standard:11111111-1111-4111-8111-111111111111",
			Name:   "Classic",
		},
		"host-1",
	)
	defer m.Dispose()

	api := New(system)
	req := httptest.NewRequest(http.MethodGet, "/api/matches", nil)
	req.Header.Set("Authorization", "Bearer test-secret")
	res := httptest.NewRecorder()

	api.getMatchesHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, res.Code, res.Body.String())
	}

	responseBody := res.Body.Bytes()
	if strings.Contains(string(responseBody), "secret-host-card") ||
		strings.Contains(string(responseBody), "secret-guest-card") {
		t.Fatal("response exposed private deck data")
	}

	var body matchesResponse
	if err := json.NewDecoder(bytes.NewReader(responseBody)).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(body.Matches) != 1 {
		t.Fatalf("expected one match, got %d", len(body.Matches))
	}

	got := body.Matches[0]
	if got.ID != m.ID || got.Name != "Private duel" {
		t.Fatalf("unexpected match identity: %#v", got)
	}
	if got.HostID != "host-1" || got.GuestID != "guest-1" ||
		got.HostUsername != "Alice" || got.GuestUsername != "Bob" {
		t.Fatalf("unexpected participants: %#v", got)
	}
	if got.Visible {
		t.Fatalf("unexpected match metadata: %#v", got)
	}
	if got.FormatID != "standard:11111111-1111-4111-8111-111111111111" ||
		got.FormatName != "Classic" {
		t.Fatalf("unexpected format descriptor: %#v", got)
	}
	if got.CreatedAt == 0 {
		t.Fatal("expected createdAt to be set")
	}
}
