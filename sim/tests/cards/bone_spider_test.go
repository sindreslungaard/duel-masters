package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const boneSpiderUID = "4d3201e8-0d9b-481e-b8e3-86cb90058e20"

// Bone Spider has the Suicide ability: it is destroyed after winning a battle.
func TestBoneSpider(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(boneSpiderUID, match.HAND)
	for range 3 {
		player.Player.SpawnCard(boneSpiderUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, boneSpiderUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	boneSpider, err := scn.FindCard(player.Player, match.BATTLEZONE, boneSpiderUID)
	require.NoError(t, err)
	assert.True(t, boneSpider.HasCondition(cnd.DestroyAfterBattle))
	assert.Equal(t, 5000, scn.Match.GetPower(boneSpider, false))
}
