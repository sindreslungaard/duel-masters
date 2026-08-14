package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const faerieChildUID = "a3cf18f0-b04f-45e9-97f7-2a2ead0a1787"

func TestFaerieChild(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(faerieChildUID, match.HAND)
	for range 4 {
		player.Player.SpawnCard(faerieChildUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, faerieChildUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	faerieChild, err := scn.FindCard(player.Player, match.BATTLEZONE, faerieChildUID)
	require.NoError(t, err)
	assert.True(t, faerieChild.HasCondition(cnd.CantBeBlocked))
	assert.Equal(t, 2000, scn.Match.GetPower(faerieChild, false))
}
