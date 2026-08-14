package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const deathbladeBeetleUID = "c1ebdda0-be88-4665-937e-2ef3ada8d378"

func TestDeathbladeBeetle(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(deathbladeBeetleUID, match.HAND)
	for range 5 {
		player.Player.SpawnCard(deathbladeBeetleUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, deathbladeBeetleUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	beetle, err := scn.FindCard(player.Player, match.BATTLEZONE, deathbladeBeetleUID)
	require.NoError(t, err)
	assert.True(t, beetle.HasCondition(cnd.DoubleBreaker))
	assert.True(t, beetle.HasCondition(cnd.PowerAttacker))
	assert.Equal(t, 3000, scn.Match.GetPower(beetle, false))
	assert.Equal(t, 7000, scn.Match.GetPower(beetle, true))
}
