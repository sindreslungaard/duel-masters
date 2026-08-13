package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const boneAssassinUID = "cc9762c3-515a-4734-a3fe-1e0c4c3b3d71"

func TestBoneAssassinTheRipper(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(boneAssassinUID, match.HAND)
	for range 4 {
		player.Player.SpawnCard(boneAssassinUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, boneAssassinUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	boneAssassin, err := scn.FindCard(player.Player, match.BATTLEZONE, boneAssassinUID)
	require.NoError(t, err)
	assert.True(t, boneAssassin.HasCondition(cnd.Slayer))
	assert.Equal(t, 2000, scn.Match.GetPower(boneAssassin, false))
}
