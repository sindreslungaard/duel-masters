package match

import (
	"bytes"
	"duel-masters/internal"
	"encoding/json"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type duelResultWebhookPlayer struct {
	UID      string `json:"uid"`
	Username string `json:"username,omitempty"`
	Deck     string `json:"deck,omitempty"`
}

type duelResultWebhookPayload struct {
	DuelID          string                   `json:"duel_id"`
	MatchName       string                   `json:"match_name,omitempty"`
	Format          string                   `json:"format"`
	StartedAt       int64                    `json:"started_at"`
	EndedAt         int64                    `json:"ended_at"`
	DurationSeconds int64                    `json:"duration_seconds"`
	Turns           int                      `json:"turns"`
	WonByDisconnect bool                     `json:"won_by_disconnect"`
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
		DuelID:          duel.UID,
		MatchName:       m.MatchName,
		Format:          duel.Format,
		StartedAt:       duel.Started,
		EndedAt:         duel.Ended,
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
		UID:      player.UID,
		Username: player.Username,
		Deck:     player.DeckStr,
	}
}

func winnerAndLoserForWebhook(winnerUID string, host *duelResultWebhookPlayer, guest *duelResultWebhookPlayer) (*duelResultWebhookPlayer, *duelResultWebhookPlayer) {
	if winnerUID == "" {
		return nil, nil
	}

	if host != nil && host.UID == winnerUID {
		return host, guest
	}

	if guest != nil && guest.UID == winnerUID {
		return guest, host
	}

	return nil, nil
}

func sendDuelResultWebhook(webhookURL string, auth string, payload duelResultWebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return &webhookStatusError{statusCode: resp.StatusCode}
	}

	return nil
}

type webhookStatusError struct {
	statusCode int
}

func (e *webhookStatusError) Error() string {
	return http.StatusText(e.statusCode)
}
