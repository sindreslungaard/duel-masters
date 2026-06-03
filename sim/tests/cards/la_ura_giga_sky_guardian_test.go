package cards

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const laUraGigaUID = "c05fe45d-690e-4856-bddb-5f46154e57e5"

func TestLaUraGigaSkyGuardian(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	player.Player.SpawnCard(laUraGigaUID, match.HAND)
	player.Player.SpawnCard(laUraGigaUID, match.MANAZONE)

	card, err := scn.FindCard(player.Player, match.HAND, laUraGigaUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	require.NoError(t, scn.ActionEndTurn(player))
	require.NoError(t, scn.ActionEndTurn(opponent))

	laUraGiga, err := scn.FindCard(player.Player, match.BATTLEZONE, laUraGigaUID)
	require.NoError(t, err)
	assert.True(t, laUraGiga.HasCondition(cnd.Blocker))
	assert.True(t, laUraGiga.HasCondition(cnd.CantAttackPlayers))
	assert.Equal(t, 2000, scn.Match.GetPower(laUraGiga, false))
}
