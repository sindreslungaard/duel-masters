package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const diaNorkUID = "f7dc24d2-2a84-46ff-9661-0b8418d68650"

func TestDiaNorkMoonlightGuardian(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(diaNorkUID, match.HAND)
	for range 4 {
		player.Player.SpawnCard(diaNorkUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, diaNorkUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	diaNork, err := scn.FindCard(player.Player, match.BATTLEZONE, diaNorkUID)
	require.NoError(t, err)
	assert.True(t, diaNork.HasCondition(cnd.Blocker))
	assert.True(t, diaNork.HasCondition(cnd.CantAttackPlayers))
	assert.Equal(t, 5000, scn.Match.GetPower(diaNork, false))
}
