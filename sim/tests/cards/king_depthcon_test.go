package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const kingDepthconUID = "cd13f7c2-aa5e-43b8-8811-700f230a5de5"

func TestKingDepthcon(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(kingDepthconUID, match.HAND)
	for range 7 {
		player.Player.SpawnCard(kingDepthconUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, kingDepthconUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	kingDepthcon, err := scn.FindCard(player.Player, match.BATTLEZONE, kingDepthconUID)
	require.NoError(t, err)
	assert.True(t, kingDepthcon.HasCondition(cnd.CantBeBlocked))
	assert.True(t, kingDepthcon.HasCondition(cnd.DoubleBreaker))
	assert.Equal(t, 6000, scn.Match.GetPower(kingDepthcon, false))
}
