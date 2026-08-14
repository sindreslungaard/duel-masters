package api

import (
	"duel-masters/game/match"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testDeck(cardID string) []string {
	return []string{cardID}
}

func createMatchRequestBody(format string, hostDeck []string, guestDeck []string) string {
	body, err := json.Marshal(matchReqBody{
		HostID:        "host-1",
		HostUsername:  "Alice",
		HostDeck:      hostDeck,
		GuestID:       "guest-1",
		GuestUsername: "Bob",
		GuestDeck:     guestDeck,
		Name:          "Authenticated duel",
		Visibility:    "public",
		Format:        format,
	})
	if err != nil {
		panic(err)
	}
	return string(body)
}

func validCreateMatchRequestBody(format string) string {
	return createMatchRequestBody(format, testDeck("host-card"), testDeck("guest-card"))
}

func TestCreateMatchHandlerRequiresServerToken(t *testing.T) {
	t.Setenv("secret", "test-secret")

	system := match.NewSystem()
	api := New(system)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/match",
		strings.NewReader(validCreateMatchRequestBody("regular")),
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
		strings.NewReader(validCreateMatchRequestBody("regular")),
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
			body: strings.Replace(validCreateMatchRequestBody("regular"), `"hostUsername":"Alice"`, `"hostUsername":" "`, 1),
		},
		{
			name: "guest username",
			body: strings.Replace(validCreateMatchRequestBody("regular"), `"guestUsername":"Bob"`, `"guestUsername":" "`, 1),
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

func TestCreateMatchHandlerRequiresDeckFieldsForEveryFormat(t *testing.T) {
	t.Setenv("secret", "test-secret")

	for _, format := range []string{"regular", "random"} {
		for _, test := range []struct {
			name      string
			hostDeck  []string
			guestDeck []string
		}{
			{name: "missing host deck", guestDeck: testDeck("guest-card")},
			{name: "missing guest deck", hostDeck: testDeck("host-card")},
		} {
			t.Run(format+"/"+test.name, func(t *testing.T) {
				system := match.NewSystem()
				api := New(system)
				req := httptest.NewRequest(
					http.MethodPost,
					"/api/match",
					strings.NewReader(createMatchRequestBody(format, test.hostDeck, test.guestDeck)),
				)
				req.Header.Set("Authorization", "Bearer test-secret")
				res := httptest.NewRecorder()

				api.createMatchHandler(res, req)

				if res.Code != http.StatusBadRequest {
					t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, res.Code, res.Body.String())
				}
				if len(system.Matches.Iter()) != 0 {
					t.Fatal("request without two valid decks created a match")
				}
			})
		}
	}
}

func TestCreateMatchHandlerAcceptsRandomFormatWithDecks(t *testing.T) {
	t.Setenv("secret", "test-secret")

	system := match.NewSystem()
	api := New(system)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/match",
		strings.NewReader(validCreateMatchRequestBody("random")),
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
	if len(matches[0].HostDeck) != 1 || len(matches[0].GuestDeck) != 1 {
		t.Fatal("created random-format match did not retain both supplied decks")
	}
	matches[0].Dispose()
}

func TestCreateMatchHandlerAcceptsEmptyDeckArrays(t *testing.T) {
	t.Setenv("secret", "test-secret")

	system := match.NewSystem()
	api := New(system)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/match",
		strings.NewReader(createMatchRequestBody("random", []string{}, []string{})),
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
	matches[0].Dispose()
}
