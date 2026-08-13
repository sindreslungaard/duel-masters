package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const zagaanUID = "07a0115e-797a-49d8-90bf-9ea6de39978d"

func TestZagaanKnightOfDarkness(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(zagaanUID, match.HAND)
	for range 6 {
		player.Player.SpawnCard(zagaanUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, zagaanUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	zagaan, err := scn.FindCard(player.Player, match.BATTLEZONE, zagaanUID)
	require.NoError(t, err)
	assert.True(t, zagaan.HasCondition(cnd.DoubleBreaker))
	assert.Equal(t, 7000, scn.Match.GetPower(zagaan, false))
}
