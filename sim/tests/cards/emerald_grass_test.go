package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const emeraldGrassUID = "ecd1ae69-4f63-4e8d-a3f4-9a5c81f98a20"

func TestEmeraldGrass(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.SpawnCard(emeraldGrassUID, match.HAND)
	player.Player.SpawnCard(emeraldGrassUID, match.MANAZONE)
	player.Player.SpawnCard(emeraldGrassUID, match.MANAZONE)

	emeraldGrassInHand, err := scn.FindCard(player.Player, match.HAND, emeraldGrassUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, emeraldGrassInHand.ID))

	emeraldGrass, err := scn.FindCard(player.Player, match.BATTLEZONE, emeraldGrassUID)
	require.NoError(t, err)
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))
	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	assert.True(t, emeraldGrass.HasCondition(cnd.Blocker))
	assert.True(t, emeraldGrass.HasCondition(cnd.CantAttackPlayers))
	assert.Equal(t, 3000, scn.Match.GetPower(emeraldGrass, false))
}
