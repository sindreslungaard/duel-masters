package cards

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const maskedHorrorUID = "ea878730-fde0-4bd0-ad25-95e49f54a1b2"

func TestMaskedHorrorShadowOfScorn(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	opponent := scn.Match.Opponent(player.Player)
	handBefore, err := opponent.Container(match.HAND)
	require.NoError(t, err)

	player.Player.SpawnCard(maskedHorrorUID, match.HAND)
	for range 5 {
		player.Player.SpawnCard(maskedHorrorUID, match.MANAZONE)
	}

	card, err := scn.FindCard(player.Player, match.HAND, maskedHorrorUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	_, err = scn.FindCard(player.Player, match.BATTLEZONE, maskedHorrorUID)
	require.NoError(t, err)

	handAfter, err := opponent.Container(match.HAND)
	require.NoError(t, err)
	assert.Len(t, handAfter, len(handBefore)-1)
}
