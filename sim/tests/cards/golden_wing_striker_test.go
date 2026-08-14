package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const goldenWingStrikerUID = "6663848d-035e-44b6-9d9f-7b236ea5bc43"

func TestGoldenWingStriker(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.SpawnCard(goldenWingStrikerUID, match.HAND)
	player.Player.SpawnCard(goldenWingStrikerUID, match.MANAZONE)
	player.Player.SpawnCard(goldenWingStrikerUID, match.MANAZONE)
	player.Player.SpawnCard(goldenWingStrikerUID, match.MANAZONE)

	goldenWingStrikerInHand, err := scn.FindCard(player.Player, match.HAND, goldenWingStrikerUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, goldenWingStrikerInHand.ID))

	goldenWingStriker, err := scn.FindCard(player.Player, match.BATTLEZONE, goldenWingStrikerUID)
	require.NoError(t, err)
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))
	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	assert.Equal(t, 2000, scn.Match.GetPower(goldenWingStriker, false))
	assert.Equal(t, 4000, scn.Match.GetPower(goldenWingStriker, true))
	assert.True(t, goldenWingStriker.HasCondition(cnd.PowerAttacker))
}
