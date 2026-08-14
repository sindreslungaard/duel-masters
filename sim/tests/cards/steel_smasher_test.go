package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const steelSmasherUID = "10e0e90f-ad7d-4b69-98d5-f01525eb1cdd"

func TestSteelSmasher(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.SpawnCard(steelSmasherUID, match.HAND)
	player.Player.SpawnCard(steelSmasherUID, match.MANAZONE)
	player.Player.SpawnCard(steelSmasherUID, match.MANAZONE)

	steelSmasherInHand, err := scn.FindCard(player.Player, match.HAND, steelSmasherUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, steelSmasherInHand.ID))

	steelSmasher, err := scn.FindCard(player.Player, match.BATTLEZONE, steelSmasherUID)
	require.NoError(t, err)
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))
	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	assert.True(t, steelSmasher.HasCondition(cnd.Creature))
	assert.True(t, steelSmasher.HasCondition(cnd.CantAttackPlayers))
	assert.False(t, steelSmasher.HasCondition(cnd.Blocker))
	assert.Equal(t, 3000, scn.Match.GetPower(steelSmasher, false))
}
