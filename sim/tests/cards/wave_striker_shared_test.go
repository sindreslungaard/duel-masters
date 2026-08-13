package cards

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"
)

// Making up the wave striker count. Their abilities are only switched on while
// 2 or more *other* creatures in the battle zone have the keyword, so every
// behavioural test needs two spare wave strikers on the board.

const (
	// Macho Melon is the cheapest wave striker with no board effect of its own,
	// which makes it the ideal filler for making up the count.
	waveStrikerFillerUID = "fa987e39-2955-4074-bcf2-b7888ae27319"
	// A small blocker and a creature comfortably above the 2000 line, for the
	// cards that care about either.
	waveStrikerSmallBlockerUID = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (blocker, 2000)
	waveStrikerBigCreatureUID  = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	waveStrikerSharedSrc       = "wave_striker_shared_test_setup"
)

// addWaveStrikerFillers puts n spare wave strikers into the player's battle
// zone. Two of them is what switches another wave striker's ability on.
func addWaveStrikerFillers(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, n int) {
	t.Helper()

	for range n {
		putCardInBattlezone(t, scn, player.Player, waveStrikerFillerUID, waveStrikerSharedSrc)
	}
}
