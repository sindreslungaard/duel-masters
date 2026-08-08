package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const kingCoralUID = "3e2940f4-5654-4456-bfc2-fa5e43911cfb"

func TestKingCoral(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(kingCoralUID, match.HAND)
	for range 3 {
		player.Player.SpawnCard(kingCoralUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, kingCoralUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	kingCoral, err := scn.FindCard(player.Player, match.BATTLEZONE, kingCoralUID)
	require.NoError(t, err)
	assert.True(t, kingCoral.HasCondition(cnd.Blocker))
	assert.Equal(t, 1000, scn.Match.GetPower(kingCoral, false))
}
