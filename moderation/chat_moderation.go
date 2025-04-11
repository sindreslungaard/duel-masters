package moderation

import (
	"bytes"
	"duel-masters/internal"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

const WriteCapacity = 5
const FlagsTolerance = 3

const (
	Hate                  = 1
	Harassment            = 1
	HarassmentThreatening = 2
	SelfHarmInstructions  = 2
	SexualMinors          = 3
)

var ChatModeration = NewChatModerationService(os.Getenv("chat_moderation_service"), os.Getenv("chat_moderation_service_secret"))

type ChatModerationService struct {
	sync.RWMutex
	enabled  bool
	endpoint string
	auth     string
	content  map[string][]string
	flags    map[string]int
	writes   int
}

func NewChatModerationService(endpoint string, auth string) *ChatModerationService {
	enabled := false

	if endpoint != "" && auth != "" {
		enabled = true
	}

	return &ChatModerationService{
		enabled:  enabled,
		endpoint: endpoint,
		auth:     auth,
		content:  make(map[string][]string),
		flags:    make(map[string]int),
		writes:   0,
	}
}

func (svc *ChatModerationService) Write(actor string, content string) {
	svc.Lock()
	defer svc.Unlock()

	s, ok := svc.content[actor]

	if ok {
		svc.content[actor] = append(s, content)
	} else {
		svc.content[actor] = []string{content}
	}

	svc.writes++

	if svc.writes >= WriteCapacity {
		svc.writes = 0
		svc.flush()
	}
}

func (svc *ChatModerationService) CheckFlags(actor string) int {
	svc.RLock()
	defer svc.RUnlock()

	flags, ok := svc.flags[actor]

	if !ok {
		return 0
	}

	return flags
}

func (svc *ChatModerationService) flush() {
	defer internal.Recover()

	logrus.Debug("Flushing chat moderation cache")

	actors := []string{}
	inputs := []string{}

	for a, s := range svc.content {
		for _, input := range s {
			actors = append(actors, a)
			inputs = append(inputs, input)
		}
	}

	for k := range svc.content {
		delete(svc.content, k)
	}

	if !svc.enabled {
		return
	}

	go func() {
		defer internal.Recover()

		logrus.Debug("Sending chatlogs to moderation service")

		payload, err := json.Marshal(map[string]any{
			"input": inputs,
		})

		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", svc.endpoint, bytes.NewBuffer(payload))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", svc.auth)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		type ModerationResponse struct {
			Results []struct {
				Flagged    bool `json:"flagged"`
				Categories struct {
					Sexual                bool `json:"sexual"`
					Hate                  bool `json:"hate"`
					Harassment            bool `json:"harassment"`
					SelfHarm              bool `json:"self-harm"`
					SexualMinors          bool `json:"sexual/minors"`
					HateThreatening       bool `json:"hate/threatening"`
					ViolenceGraphic       bool `json:"violence/graphic"`
					SelfHarmIntent        bool `json:"self-harm/intent"`
					SelfHarmInstructions  bool `json:"self-harm/instructions"`
					HarassmentThreatening bool `json:"harassment/threatening"`
					Violence              bool `json:"violence"`
				} `json:"categories"`
			} `json:"results"`
		}

		var moderationResp ModerationResponse
		if err := json.Unmarshal(body, &moderationResp); err != nil {
			panic(err)
		}

		svc.Lock()
		defer svc.Unlock()

		for i, result := range moderationResp.Results {
			actor := actors[i]
			points := 0

			if result.Categories.Hate {
				points += Hate
			}

			if result.Categories.Harassment {
				points += Harassment
			}

			if result.Categories.HarassmentThreatening {
				points += HarassmentThreatening
			}

			if result.Categories.SelfHarmInstructions {
				points += SelfHarmInstructions
			}

			if result.Categories.SexualMinors {
				points += SexualMinors
			}

			if points > 0 {
				p, ok := svc.flags[actor]

				if ok {
					svc.flags[actor] = p + points
				} else {
					svc.flags[actor] = points
				}
			}
		}

	}()
}

func (svc *ChatModerationService) ResetFlags() {
	svc.Lock()
	defer svc.Unlock()

	for k := range svc.flags {
		delete(svc.flags, k)
	}
}

func (svc *ChatModerationService) Toggle() bool {
	if svc.enabled {
		svc.enabled = false
	} else {
		svc.enabled = true
	}

	svc.ResetFlags()

	return svc.enabled
}

func (svc *ChatModerationService) Flags() map[string]int {
	svc.RLock()
	defer svc.RUnlock()

	m := make(map[string]int)

	for k, v := range svc.flags {
		m[k] = v
	}

	return m
}
