package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const deathligerUID = "dc1b51b3-52e7-4f1c-8770-515d4e1cb53d"

func TestDeathligerLionOfChaos(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(deathligerUID, match.HAND)
	for range 7 {
		player.Player.SpawnCard(deathligerUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, deathligerUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	deathliger, err := scn.FindCard(player.Player, match.BATTLEZONE, deathligerUID)
	require.NoError(t, err)
	assert.True(t, deathliger.HasCondition(cnd.DoubleBreaker))
	assert.Equal(t, 9000, scn.Match.GetPower(deathliger, false))
}
