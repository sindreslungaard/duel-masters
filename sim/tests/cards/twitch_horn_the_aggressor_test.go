package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const twitchHornTheAggressorUID = "47875b7c-6472-41d9-8994-7c21306a1a99"

func TestTwitchHornTheAggressor(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer().Player

	player.SpawnCard(twitchHornTheAggressorUID, match.HAND)
	twitchHorn, err := scn.FindCard(player, match.HAND, twitchHornTheAggressorUID)
	require.NoError(t, err)
	moved, err := player.MoveCard(twitchHorn.ID, match.HAND, match.BATTLEZONE, "twitch_horn_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)

	for range 3 {
		player.SpawnCard(scowlingTomatoUID, match.MANAZONE)
	}
	mana, err := player.Container(match.MANAZONE)
	require.NoError(t, err)
	require.Len(t, mana, 3)
	mana[0].Tapped = true
	mana[1].Tapped = true

	assert.Equal(t, "Twitch Horn, the Aggressor", twitchHorn.Name)
	assert.Equal(t, 2000, twitchHorn.Power)
	assert.Equal(t, 6, twitchHorn.ManaCost)
	assert.Equal(t, civ.Nature, twitchHorn.Civ)
	assert.True(t, twitchHorn.HasFamily(family.HornedBeast))
	assert.Equal(t, 2000, scn.Match.GetPower(twitchHorn, false))
	assert.Equal(t, 6000, scn.Match.GetPower(twitchHorn, true))

	_, err = player.MoveCard(twitchHorn.ID, match.BATTLEZONE, match.HAND, "twitch_horn_test_setup")
	require.NoError(t, err)
	assert.Equal(t, 2000, scn.Match.GetPower(twitchHorn, true), "the ability is inactive outside the battle zone")
}
