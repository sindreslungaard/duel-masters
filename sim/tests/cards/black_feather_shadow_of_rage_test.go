package cards

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const blackFeatherUID = "c62208ec-efc8-4b08-bb01-6cd5251b0969"

// Black Feather, Shadow of Rage destroys one of the controller's creatures when summoned.
// With no other creature in the battlezone it is forced to destroy itself.
func TestBlackFeatherShadowOfRage_DestroysItself(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.SpawnCard(blackFeatherUID, match.HAND)
	player.Player.SpawnCard(blackFeatherUID, match.MANAZONE)

	card, err := scn.FindCard(player.Player, match.HAND, blackFeatherUID)
	require.NoError(t, err)
	require.NoError(t, scn.ActionPlayCard(player, card.ID))

	assert.Eventually(t, func() bool {
		_, err := scn.FindCard(player.Player, match.GRAVEYARD, blackFeatherUID)
		return err == nil
	}, time.Second, 10*time.Millisecond)
}
