package tests

import (
	"duel-masters/tests/scenario"
	"fmt"
	"testing"
)

// No reason to run this in CI - change below const to true before running locally.
// Run the test in verbose mode to see the printed result.
const ShufflerTestEnabled = false
const ShufflerTestIterations = 10000

func TestShuffler(t *testing.T) {
	scenario := scenario.New(scenario.Options{})

	t.Run("Shuffler shuffles as expected", func(t *testing.T) {
		if !ShufflerTestEnabled {
			t.SkipNow()
		}

		fmt.Println("Running shuffler test...")

		// Result holds the number of times each card appears in the first 5 cards after shuffling
		result := map[string]int{}

		p := scenario.Match.Player1.Player

		for i := 0; i < ShufflerTestIterations; i++ {
			p.ShuffleDeck()
			hand := p.PeekDeck(5)

			for _, card := range hand {
				result[card.ID]++
			}
		}

		for k, v := range result {
			fmt.Printf("%.1f%% Card ID: %s\n", (float64(v)/float64(ShufflerTestIterations))*100, k)
		}
	})
}
