package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	gigarayzeUID       = "5ead68c2-dce0-4e8e-b718-02f2cce5dfa0"
	gigarayzeWaterUID  = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (water)
	gigarayzeNatureUID = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur (nature)
	gigarayzeSetupSrc  = "gigarayze_test_setup"
)

func TestGigarayze(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// Kept in hand: it may ask for a target as it arrives, and a card moved
		// into play from the test goroutine would leave it waiting on a prompt
		// only it could answer.
		card := spawnForLater(t, player, gigarayzeUID)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Gigarayze", 2000, 4, []string{civ.Darkness})
		assert.True(t, card.HasFamily(family.Chimera))
	})

	t.Run("it may take a water or fire creature back", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		water, err := player.Player.SpawnCard(gigarayzeWaterUID, match.GRAVEYARD)
		require.NoError(t, err)
		nature, err := player.Player.SpawnCard(gigarayzeNatureUID, match.GRAVEYARD)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		summonWithOwnMana(t, scn, player, gigarayzeUID)

		answerInTurn(t, scn, player, water.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, water.Zone)
		assert.Equal(t, match.GRAVEYARD, nature.Zone, "nature is not water or fire")
	})

	t.Run("the return may be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		water, err := player.Player.SpawnCard(gigarayzeWaterUID, match.GRAVEYARD)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		summonWithOwnMana(t, scn, player, gigarayzeUID)

		cancelInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, water.Zone)
	})
}
