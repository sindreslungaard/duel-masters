package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	torpedoClusterUID     = "5c343af6-6d86-4b1c-a5e7-d847650f4038"
	torpedoClusterManaUID = "9781089f-1aa9-4a75-b106-35e9d431e31d" // Aqua Vehicle
)

func TestTorpedoCluster(t *testing.T) {
	t.Run("returns a chosen card from its owner's mana zone to their hand", func(t *testing.T) {
		scn, player, torpedo := setupTorpedoClusterTest(t)

		assert.Equal(t, "Torpedo Cluster", torpedo.Name)
		assert.Equal(t, 3000, torpedo.Power)
		assert.Equal(t, 3, torpedo.ManaCost)
		assert.Equal(t, civ.Water, torpedo.Civ)
		assert.True(t, torpedo.HasFamily(family.CyberCluster))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, torpedo.ID))
		assert.Equal(t, match.BATTLEZONE, torpedo.Zone)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, action.MinSelections)
		assert.Equal(t, 1, action.MaxSelections)
		assert.False(t, action.Cancellable, "the effect is mandatory")

		manaCards, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		require.Len(t, manaCards, 3)
		chosen := manaCards[0]

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, chosen.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.HAND, chosen.Zone)
		remaining, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, 2)
	})

	t.Run("a rejected cancellation still resolves the mandatory selection", func(t *testing.T) {
		scn, player, torpedo := setupTorpedoClusterTest(t)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, torpedo.ID))

		_, err = scn.LatestAction(player, promptStart)
		require.NoError(t, err)

		warningStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForMessage(player, warningStart, "action_error"))

		manaCards, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		require.Len(t, manaCards, 3, "cancelling must not move anything")
		chosen := manaCards[0]

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, chosen.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.HAND, chosen.Zone)
	})
}

func setupTorpedoClusterTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.Card) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.SpawnCard(torpedoClusterUID, match.HAND)
	for range 3 {
		player.Player.SpawnCard(torpedoClusterManaUID, match.MANAZONE)
	}

	torpedo, err := scn.FindCard(player.Player, match.HAND, torpedoClusterUID)
	require.NoError(t, err)

	return scn, player, torpedo
}
