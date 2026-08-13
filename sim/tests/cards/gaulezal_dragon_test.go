package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const gaulezalDragonUID = "79c48731-193b-4dc6-b26f-1eb820357367"

func TestGaulezalDragon(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.SpawnCard(gaulezalDragonUID, match.HAND)
	dragon, err := scn.FindCard(player.Player, match.HAND, gaulezalDragonUID)
	require.NoError(t, err)

	assert.Equal(t, "Gaulezal Dragon", dragon.Name)
	assert.Equal(t, 11000, dragon.Power)
	assert.Equal(t, 9, dragon.ManaCost)
	assert.Equal(t, []string{civ.Fire}, dragon.Civs)
	assert.True(t, dragon.HasFamily(family.ArmoredDragon))

	moved, err := player.Player.MoveCard(dragon.ID, match.HAND, match.BATTLEZONE, "gaulezal_dragon_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	require.NoError(t, scn.ActionEndTurn(player))

	assert.True(t, dragon.HasCondition(cnd.DoubleBreaker))
}
