package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const granGureUID = "39090f65-779c-46c9-856c-67303dd5605c"

func TestGranGureSpaceGuardian(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(granGureUID, match.HAND)
	for range 6 {
		player.Player.SpawnCard(granGureUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, granGureUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	granGure, err := scn.FindCard(player.Player, match.BATTLEZONE, granGureUID)
	require.NoError(t, err)
	assert.True(t, granGure.HasCondition(cnd.Blocker))
	assert.True(t, granGure.HasCondition(cnd.CantAttackPlayers))
	assert.Equal(t, 9000, scn.Match.GetPower(granGure, false))
}
