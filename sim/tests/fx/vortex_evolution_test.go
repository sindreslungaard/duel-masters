package fx

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Wise Starnoid, Avatar of Hope needs a Light Bringer and a Cyber Lord.
	vortexEvolutionUID             = "ba894f7e-b7d7-409e-8393-cab4285a879c"
	vortexEvolutionLightBringerUID = "7b58e8c2-0b1e-4ef5-812f-e667c2092c73" // Reusol, the Oracle
	vortexEvolutionCyberLordUID    = "7a6f1c82-a8ac-4646-b3e9-fb8592bdd0a4" // Tropico
)

// summonVortexEvolution pays for uid with copies of itself and plays it, then
// answers the stacking-order prompt that appears once its controller has
// exactly one of each required creature: with no real choice for either
// requirement, both selections resolve on their own.
func summonVortexEvolution(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, uid string, topID string, bottomID string) *match.Card {
	t.Helper()

	card, err := player.Player.SpawnCard(uid, match.HAND)
	require.NoError(t, err)

	for range card.ManaCost {
		_, err := player.Player.SpawnCard(uid, match.MANAZONE)
		require.NoError(t, err)
	}

	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	action, err := scn.LatestAction(player, 0)
	require.NoError(t, err, "expected the stacking-order prompt")
	require.Len(t, action.Cards, 2)

	require.NoError(t, scn.SubmitAction(player, topID, bottomID))
	require.NoError(t, scn.WaitForEventLoop())

	return card
}

// TestVortexEvolutionDispatchesEvolutionEventForBothBases guards a Vortex
// Evolution's compatibility with anything that generically watches for "a
// creature evolved" (fx.CreatureEvolved) rather than a specific card: with two
// bases instead of one, both must be announced.
func TestVortexEvolutionDispatchesEvolutionEventForBothBases(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	watcher := &evolutionEventWatcher{}
	observer, err := player.Player.SpawnCard(evolutionEventBaseUID, match.HAND)
	require.NoError(t, err)
	watcher.attach(observer)

	lightBringer := putInBattlezone(t, scn, player.Player, vortexEvolutionLightBringerUID)
	cyberLord := putInBattlezone(t, scn, player.Player, vortexEvolutionCyberLordUID)

	evolution := summonVortexEvolution(t, scn, player, vortexEvolutionUID, lightBringer.ID, cyberLord.ID)

	require.Len(t, watcher.events, 2)
	assert.Equal(t, evolution.ID, watcher.events[0].CardID)
	assert.Equal(t, evolution.ID, watcher.events[1].CardID)
	assert.ElementsMatch(t, []string{lightBringer.ID, cyberLord.ID}, []string{watcher.events[0].BaseID, watcher.events[1].BaseID})
}

// TestVortexEvolutionEntersTappedWhenItsTopBaseIsTapped is the vortex
// counterpart of TestEvolutionEntersTappedWhenItsBaseIsTapped: it must not
// regress the same way a plain evolution did.
func TestVortexEvolutionEntersTappedWhenItsTopBaseIsTapped(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	lightBringer := putInBattlezone(t, scn, player.Player, vortexEvolutionLightBringerUID)
	cyberLord := putInBattlezone(t, scn, player.Player, vortexEvolutionCyberLordUID)
	cyberLord.Tapped = true

	evolution := summonVortexEvolution(t, scn, player, vortexEvolutionUID, cyberLord.ID, lightBringer.ID)

	require.Equal(t, match.BATTLEZONE, evolution.Zone)
	assert.True(t, evolution.Tapped, "it took the tap state of the base it was put on top of")
}
