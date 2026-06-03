package cards

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const aquaHulcusUID = "57eeb3c3-2561-4841-a381-2e50d17533d1"
const aquaHulcusDeckSeedUID = "c5a869f4-a959-4667-a352-92df5369e0b9"

func TestAquaHulcus(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.DestroyDeck()
	player.Player.SpawnCard(aquaHulcusDeckSeedUID, match.DECK)
	player.Player.SpawnCard(aquaHulcusUID, match.DECK)

	player.Player.SpawnCard(aquaHulcusUID, match.HAND)
	player.Player.SpawnCard(aquaHulcusUID, match.MANAZONE)
	player.Player.SpawnCard(aquaHulcusUID, match.MANAZONE)
	player.Player.SpawnCard(aquaHulcusUID, match.MANAZONE)

	handBefore, err := player.Player.Container(match.HAND)
	require.NoError(t, err)
	deckBefore, err := player.Player.Container(match.DECK)
	require.NoError(t, err)

	aquaHulcus, err := scn.FindCard(player.Player, match.HAND, aquaHulcusUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, aquaHulcus.ID))
	require.NoError(t, scn.SubmitAction(player))

	assert.Eventually(t, func() bool {
		handAfter, err := player.Player.Container(match.HAND)
		if err != nil {
			return false
		}

		deckAfter, err := player.Player.Container(match.DECK)
		if err != nil {
			return false
		}

		if _, err := scn.FindCard(player.Player, match.HAND, aquaHulcusDeckSeedUID); err != nil {
			return false
		}

		return len(handAfter) == len(handBefore) && len(deckAfter) == len(deckBefore)-1
	}, time.Second, 10*time.Millisecond)
}
