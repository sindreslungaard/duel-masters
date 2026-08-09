package api

import (
	"duel-masters/game/match"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const createMatchRequestBody = `{
	"hostId":"host-1",
	"hostUsername":"Alice",
	"hostDeck":[],
	"guestId":"guest-1",
	"guestUsername":"Bob",
	"guestDeck":[],
	"name":"Authenticated duel",
	"visibility":"public",
	"format":"regular"
}`

func TestCreateMatchHandlerRequiresServerToken(t *testing.T) {
	t.Setenv("secret", "test-secret")

	system := match.NewSystem()
	api := New(system)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/match",
		strings.NewReader(createMatchRequestBody),
	)
	res := httptest.NewRecorder()

	api.createMatchHandler(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
	if len(system.Matches.Iter()) != 0 {
		t.Fatal("unauthorized request created a match")
	}
}

func TestCreateMatchHandlerAcceptsServerToken(t *testing.T) {
	t.Setenv("secret", "test-secret")

	system := match.NewSystem()
	api := New(system)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/match",
		strings.NewReader(createMatchRequestBody),
	)
	req.Header.Set("Authorization", "Bearer test-secret")
	res := httptest.NewRecorder()

	api.createMatchHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, res.Code, res.Body.String())
	}

	matches := system.Matches.Iter()
	if len(matches) != 1 {
		t.Fatalf("expected one match, got %d", len(matches))
	}
	summary := matches[0].Summary()
	if summary.HostUsername != "Alice" || summary.GuestUsername != "Bob" {
		t.Fatalf("unexpected creation-time usernames: %#v", summary)
	}
	matches[0].Dispose()
}

func TestCreateMatchHandlerRequiresUsernames(t *testing.T) {
	t.Setenv("secret", "test-secret")

	for _, test := range []struct {
		name string
		body string
	}{
		{
			name: "host username",
			body: strings.Replace(createMatchRequestBody, `"hostUsername":"Alice"`, `"hostUsername":" "`, 1),
		},
		{
			name: "guest username",
			body: strings.Replace(createMatchRequestBody, `"guestUsername":"Bob"`, `"guestUsername":" "`, 1),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			system := match.NewSystem()
			api := New(system)
			req := httptest.NewRequest(
				http.MethodPost,
				"/api/match",
				strings.NewReader(test.body),
			)
			req.Header.Set("Authorization", "Bearer test-secret")
			res := httptest.NewRecorder()

			api.createMatchHandler(res, req)

			if res.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
			}
			if len(system.Matches.Iter()) != 0 {
				t.Fatal("request without both usernames created a match")
			}
		})
	}
}
