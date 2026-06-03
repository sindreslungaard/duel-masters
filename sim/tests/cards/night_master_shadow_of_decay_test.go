package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const nightMasterUID = "f16795cc-4378-4e36-b13a-19f9b932228c"

func TestNightMasterShadowOfDecay(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(nightMasterUID, match.HAND)
	for range 6 {
		player.Player.SpawnCard(nightMasterUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, nightMasterUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	nightMaster, err := scn.FindCard(player.Player, match.BATTLEZONE, nightMasterUID)
	require.NoError(t, err)
	assert.True(t, nightMaster.HasCondition(cnd.Blocker))
	assert.Equal(t, 3000, scn.Match.GetPower(nightMaster, false))
}
