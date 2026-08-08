package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const darkRavenUID = "162f70fb-33f7-4436-a114-41f255c0ce7e"

func TestDarkRavenShadowOfGrief(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(darkRavenUID, match.HAND)
	for range 4 {
		player.Player.SpawnCard(darkRavenUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, darkRavenUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	darkRaven, err := scn.FindCard(player.Player, match.BATTLEZONE, darkRavenUID)
	require.NoError(t, err)
	assert.True(t, darkRaven.HasCondition(cnd.Blocker))
	assert.Equal(t, 1000, scn.Match.GetPower(darkRaven, false))
}
