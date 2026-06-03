package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const wanderingBraineaterUID = "90b2ed59-828c-4237-ac2e-b7008a02ad2e"

func TestWanderingBraineater(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(wanderingBraineaterUID, match.HAND)
	for range 2 {
		player.Player.SpawnCard(wanderingBraineaterUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, wanderingBraineaterUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	braineater, err := scn.FindCard(player.Player, match.BATTLEZONE, wanderingBraineaterUID)
	require.NoError(t, err)
	assert.True(t, braineater.HasCondition(cnd.Blocker))
	assert.True(t, braineater.HasCondition(cnd.CantAttackPlayers))
	assert.True(t, braineater.HasCondition(cnd.CantAttackCreatures))
	assert.Equal(t, 2000, scn.Match.GetPower(braineater, false))
}
