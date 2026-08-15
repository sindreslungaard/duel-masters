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

func createMatchRequestBody(hostDeck []string, guestDeck []string) string {
	return createMatchRequestBodyWithDescriptor(hostDeck, guestDeck, "", "")
}

func createMatchRequestBodyWithDescriptor(
	hostDeck []string,
	guestDeck []string,
	formatID string,
	formatName string,
) string {
	body, err := json.Marshal(matchReqBody{
		HostID:        "host-1",
		HostUsername:  "Alice",
		HostDeck:      hostDeck,
		GuestID:       "guest-1",
		GuestUsername: "Bob",
		GuestDeck:     guestDeck,
		Name:          "Authenticated duel",
		Visibility:    "public",
		FormatID:      formatID,
		FormatName:    formatName,
	})
	if err != nil {
		panic(err)
	}
	return string(body)
}

func validCreateMatchRequestBody() string {
	return createMatchRequestBody(testDeck("host-card"), testDeck("guest-card"))
}

func TestCreateMatchHandlerRequiresServerToken(t *testing.T) {
	t.Setenv("secret", "test-secret")

	system := match.NewSystem()
	api := New(system)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/match",
		strings.NewReader(validCreateMatchRequestBody()),
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
		strings.NewReader(validCreateMatchRequestBody()),
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

func TestCreateMatchHandlerStoresFormatDescriptor(t *testing.T) {
	t.Setenv("secret", "test-secret")

	body := createMatchRequestBodyWithDescriptor(
		testDeck("host-card"),
		testDeck("guest-card"),
		"custom:22222222-2222-4222-8222-222222222222",
		"  Custom  ",
	)

	system := match.NewSystem()
	api := New(system)
	req := httptest.NewRequest(http.MethodPost, "/api/match", strings.NewReader(body))
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
	defer matches[0].Dispose()

	summary := matches[0].Summary()
	if summary.FormatID != "custom:22222222-2222-4222-8222-222222222222" {
		t.Fatalf("unexpected format id: %#v", summary)
	}
	if summary.FormatName != "Custom" {
		t.Fatalf("expected the format name to be trimmed, got %q", summary.FormatName)
	}
	if strings.Contains(res.Body.String(), "formatId") {
		t.Fatal("format descriptor leaked into the player-facing match payload")
	}
}

func TestCreateMatchHandlerAllowsMissingFormatDescriptor(t *testing.T) {
	t.Setenv("secret", "test-secret")

	system := match.NewSystem()
	api := New(system)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/match",
		strings.NewReader(validCreateMatchRequestBody()),
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
	defer matches[0].Dispose()

	summary := matches[0].Summary()
	if summary.FormatID != "" || summary.FormatName != "" {
		t.Fatalf("expected an empty format descriptor, got %#v", summary)
	}
}

func TestCreateMatchHandlerRejectsOversizedFormatDescriptor(t *testing.T) {
	t.Setenv("secret", "test-secret")

	for _, test := range []struct {
		name       string
		formatID   string
		formatName string
	}{
		{name: "format id", formatID: strings.Repeat("a", 101)},
		{name: "format name", formatName: strings.Repeat("a", 81)},
	} {
		t.Run(test.name, func(t *testing.T) {
			body := createMatchRequestBodyWithDescriptor(
				testDeck("host-card"),
				testDeck("guest-card"),
				test.formatID,
				test.formatName,
			)

			system := match.NewSystem()
			api := New(system)
			req := httptest.NewRequest(http.MethodPost, "/api/match", strings.NewReader(body))
			req.Header.Set("Authorization", "Bearer test-secret")
			res := httptest.NewRecorder()

			api.createMatchHandler(res, req)

			if res.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
			}
			if len(system.Matches.Iter()) != 0 {
				t.Fatal("an oversized format descriptor created a match")
			}
		})
	}
}

func TestCreateMatchHandlerRequiresUsernames(t *testing.T) {
	t.Setenv("secret", "test-secret")

	for _, test := range []struct {
		name string
		body string
	}{
		{
			name: "host username",
			body: strings.Replace(validCreateMatchRequestBody(), `"hostUsername":"Alice"`, `"hostUsername":" "`, 1),
		},
		{
			name: "guest username",
			body: strings.Replace(validCreateMatchRequestBody(), `"guestUsername":"Bob"`, `"guestUsername":" "`, 1),
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

func TestCreateMatchHandlerRequiresDeckFields(t *testing.T) {
	t.Setenv("secret", "test-secret")

	{
		for _, test := range []struct {
			name      string
			hostDeck  []string
			guestDeck []string
		}{
			{name: "missing host deck", guestDeck: testDeck("guest-card")},
			{name: "missing guest deck", hostDeck: testDeck("host-card")},
		} {
			t.Run(test.name, func(t *testing.T) {
				system := match.NewSystem()
				api := New(system)
				req := httptest.NewRequest(
					http.MethodPost,
					"/api/match",
					strings.NewReader(createMatchRequestBody(test.hostDeck, test.guestDeck)),
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

func TestCreateMatchHandlerRetainsSuppliedDecks(t *testing.T) {
	t.Setenv("secret", "test-secret")

	system := match.NewSystem()
	api := New(system)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/match",
		strings.NewReader(validCreateMatchRequestBody()),
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
		t.Fatal("created match did not retain both supplied decks")
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
		strings.NewReader(createMatchRequestBody([]string{}, []string{})),
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
