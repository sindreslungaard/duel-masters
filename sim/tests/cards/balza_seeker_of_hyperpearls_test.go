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

const balzaSeekerOfHyperpearlsUID = "d4d00738-81f8-4782-91e3-de96f40023d9"

func TestBalzaSeekerOfHyperpearls(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.SpawnCard(balzaSeekerOfHyperpearlsUID, match.HAND)
	balza, err := scn.FindCard(player.Player, match.HAND, balzaSeekerOfHyperpearlsUID)
	require.NoError(t, err)

	assert.Equal(t, "Balza, Seeker of Hyperpearls", balza.Name)
	assert.Equal(t, 4000, balza.Power)
	assert.Equal(t, 8, balza.ManaCost)
	assert.Equal(t, []string{civ.Light}, balza.Civs)
	assert.True(t, balza.HasFamily(family.MechaThunder))

	moved, err := player.Player.MoveCard(balza.ID, match.HAND, match.BATTLEZONE, "balza_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	require.NoError(t, scn.ActionEndTurn(player))

	assert.True(t, balza.HasCondition(cnd.ShieldTrigger))
	assert.True(t, balza.HasCondition(cnd.Blocker))
	assert.True(t, balza.HasCondition(cnd.CantAttackPlayers))
}
