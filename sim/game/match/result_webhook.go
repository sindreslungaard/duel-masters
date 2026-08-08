package match

import (
	"bytes"
	"duel-masters/internal"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type duelResultWebhookPlayer struct {
	UserID   string `json:"userId"`
	Username string `json:"username,omitempty"`
	Deck     string `json:"deck,omitempty"`
}

type duelResultWebhookPayload struct {
	MatchID         string                   `json:"matchId"`
	MatchName       string                   `json:"matchName,omitempty"`
	Format          string                   `json:"format"`
	StartedAt       string                   `json:"startedAt"`
	EndedAt         string                   `json:"endedAt"`
	DurationSeconds int64                    `json:"durationSeconds"`
	Turns           int                      `json:"turns"`
	WonByDisconnect bool                     `json:"wonByDisconnect"`
	Host            *duelResultWebhookPlayer `json:"host,omitempty"`
	Guest           *duelResultWebhookPlayer `json:"guest,omitempty"`
	Winner          *duelResultWebhookPlayer `json:"winner,omitempty"`
	Loser           *duelResultWebhookPlayer `json:"loser,omitempty"`
}

func (m *Match) sendMatchResultWebhook(duel DuelRecord) {
	webhookURL := os.Getenv("duel_result_webhook_url")
	if webhookURL == "" {
		logrus.Debug("Duel result webhook not configured, skipping")
		return
	}

	auth := os.Getenv("duel_token_secret")
	if auth == "" {
		logrus.Warn("duel_result_webhook_url is set, but duel_token_secret is missing; skipping match result webhook")
		return
	}

	payload := m.newDuelResultWebhookPayload(duel)

	go func() {
		defer internal.Recover()

		if err := sendDuelResultWebhook(webhookURL, auth, payload); err != nil {
			logrus.WithError(err).Error("Failed to send duel result webhook")
		}
	}()
}

func (m *Match) newDuelResultWebhookPayload(duel DuelRecord) duelResultWebhookPayload {
	host := newDuelResultWebhookPlayer(m.Player1)
	guest := newDuelResultWebhookPlayer(m.Player2)
	winner, loser := winnerAndLoserForWebhook(duel.Winner, host, guest)

	durationSeconds := duel.Ended - duel.Started
	if durationSeconds < 0 {
		durationSeconds = 0
	}

	return duelResultWebhookPayload{
		MatchID:         duel.UID,
		MatchName:       m.MatchName,
		Format:          duel.Format,
		StartedAt:       time.Unix(duel.Started, 0).UTC().Format(time.RFC3339),
		EndedAt:         time.Unix(duel.Ended, 0).UTC().Format(time.RFC3339),
		DurationSeconds: durationSeconds,
		Turns:           duel.Turns,
		WonByDisconnect: duel.WonByDisconnect,
		Host:            host,
		Guest:           guest,
		Winner:          winner,
		Loser:           loser,
	}
}

func newDuelResultWebhookPlayer(player *PlayerReference) *duelResultWebhookPlayer {
	if player == nil {
		return nil
	}

	return &duelResultWebhookPlayer{
		UserID:   player.UID,
		Username: player.Username,
		Deck:     player.DeckStr,
	}
}

func winnerAndLoserForWebhook(winnerUID string, host *duelResultWebhookPlayer, guest *duelResultWebhookPlayer) (*duelResultWebhookPlayer, *duelResultWebhookPlayer) {
	if winnerUID == "" {
		return nil, nil
	}

	if host != nil && host.UserID == winnerUID {
		return host, guest
	}

	if guest != nil && guest.UserID == winnerUID {
		return guest, host
	}

	return nil, nil
}

func sendDuelResultWebhook(webhookURL string, secret string, payload duelResultWebhookPayload) error {
	req, err := newDuelResultWebhookRequest(webhookURL, secret, payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return &webhookStatusError{statusCode: resp.StatusCode}
	}

	return nil
}

func newDuelResultWebhookRequest(webhookURL string, secret string, payload duelResultWebhookPayload) (*http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(secret))

	return req, nil
}

type webhookStatusError struct {
	statusCode int
}

func (e *webhookStatusError) Error() string {
	return fmt.Sprintf("duel result webhook returned HTTP %d (%s)", e.statusCode, http.StatusText(e.statusCode))
}
