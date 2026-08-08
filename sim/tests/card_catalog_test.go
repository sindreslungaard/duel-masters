package tests

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"
)

const expectedDM01ThroughDM12Cards = 890

var (
	dm01ThroughDM12Pattern = regexp.MustCompile(`^DM-(0[1-9]|1[0-2])(?:\s|$)`)
	uuidPattern            = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

func TestDM01ThroughDM12CardsHaveUniqueUUIDs(t *testing.T) {
	t.Parallel()

	type printing struct {
		Set string `json:"set"`
	}
	type card struct {
		ID        string     `json:"id"`
		Name      string     `json:"name"`
		Printings []printing `json:"printings"`
	}
	type catalog struct {
		Cards []card `json:"cards"`
	}

	data, err := os.ReadFile("../DuelMastersCards.json")
	if err != nil {
		t.Fatalf("read card catalog: %v", err)
	}

	var cards catalog
	if err := json.Unmarshal(data, &cards); err != nil {
		t.Fatalf("parse card catalog: %v", err)
	}

	seen := make(map[string]string)
	covered := 0

	for _, card := range cards.Cards {
		inCoveredSet := false
		for _, printing := range card.Printings {
			if dm01ThroughDM12Pattern.MatchString(printing.Set) {
				inCoveredSet = true
				break
			}
		}

		if !inCoveredSet {
			continue
		}

		covered++
		if !uuidPattern.MatchString(card.ID) {
			t.Errorf("card %q has invalid or missing UUID %q", card.Name, card.ID)
			continue
		}

		if previous, exists := seen[card.ID]; exists {
			t.Errorf("cards %q and %q share UUID %s", previous, card.Name, card.ID)
			continue
		}
		seen[card.ID] = card.Name
	}

	if covered != expectedDM01ThroughDM12Cards {
		t.Fatalf("expected %d unique DM-01 through DM-12 cards, found %d", expectedDM01ThroughDM12Cards, covered)
	}
}
