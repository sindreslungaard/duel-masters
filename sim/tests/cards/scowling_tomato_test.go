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

const scowlingTomatoUID = "2e10b4fb-3f85-4144-8762-51c04fe609d5"

func TestScowlingTomato(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.SpawnCard(scowlingTomatoUID, match.HAND)
	tomato, err := scn.FindCard(player.Player, match.HAND, scowlingTomatoUID)
	require.NoError(t, err)

	assert.Equal(t, "Scowling Tomato", tomato.Name)
	assert.Equal(t, 2000, tomato.Power)
	assert.Equal(t, 2, tomato.ManaCost)
	assert.Equal(t, civ.Nature, tomato.Civ)
	assert.True(t, tomato.HasFamily(family.WildVeggies))

	require.NoError(t, scn.ActionEndTurn(player))
	assert.True(t, tomato.HasCondition(cnd.Creature))
}
