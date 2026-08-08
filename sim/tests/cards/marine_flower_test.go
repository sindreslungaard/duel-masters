package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const marineFlowerUID = "3f331274-f5f8-42e7-9f28-ce637add34d4"

func TestMarineFlower(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(marineFlowerUID, match.HAND)
	player.Player.SpawnCard(marineFlowerUID, match.MANAZONE)

	card, err := scn.FindCard(player.Player, match.HAND, marineFlowerUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	marineFlower, err := scn.FindCard(player.Player, match.BATTLEZONE, marineFlowerUID)
	require.NoError(t, err)
	assert.True(t, marineFlower.HasCondition(cnd.Blocker))
	assert.True(t, marineFlower.HasCondition(cnd.CantAttackPlayers))
	assert.True(t, marineFlower.HasCondition(cnd.CantAttackCreatures))
	assert.Equal(t, 2000, scn.Match.GetPower(marineFlower, false))
}
