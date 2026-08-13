package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	diamondiaUID          = "6a668bef-98c8-4b55-9412-414f25e914e9"
	diamondiaSnowFaerieID = "41b82092-4d34-4b3c-934a-d7bd4fca654a" // Garabon, the Glider (Snow Faerie)
	diamondiaSetupSrc     = "diamondia_the_blizzard_rider_test_setup"
)

func TestDiamondiaTheBlizzardRider(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		diamondia := putCardInBattlezone(t, scn, player.Player, diamondiaUID, diamondiaSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, diamondia, "Diamondia, the Blizzard Rider", 5000, 3, []string{civ.Nature})
		assert.True(t, diamondia.HasFamily(family.SnowFaerie))
		assert.True(t, diamondia.HasCondition(cnd.Evolution))
	})

	t.Run("it evolves and gathers every Snow Faerie from the graveyard and mana", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		base := putCardInBattlezone(t, scn, player.Player, diamondiaSnowFaerieID, diamondiaSetupSrc)

		inGrave, err := player.Player.SpawnCard(diamondiaSnowFaerieID, match.GRAVEYARD)
		require.NoError(t, err)
		inMana, err := player.Player.SpawnCard(diamondiaSnowFaerieID, match.MANAZONE)
		require.NoError(t, err)
		otherInGrave, err := player.Player.SpawnCard(immortalBaronVorgUID, match.GRAVEYARD)
		require.NoError(t, err)
		theirs, err := opponent.Player.SpawnCard(diamondiaSnowFaerieID, match.GRAVEYARD)
		require.NoError(t, err)

		diamondia, err := player.Player.SpawnCard(diamondiaUID, match.HAND)
		require.NoError(t, err)
		for range 3 {
			_, err := player.Player.SpawnCard(diamondiaUID, match.MANAZONE)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionPlayCard(player, diamondia.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, diamondia.Zone)
		assert.Equal(t, match.HIDDENZONE, base.Zone, "the base goes under the evolution")
		assert.Equal(t, match.HAND, inGrave.Zone)
		assert.Equal(t, match.HAND, inMana.Zone)
		assert.Equal(t, match.GRAVEYARD, otherInGrave.Zone, "only Snow Faeries come back")
		assert.Equal(t, match.GRAVEYARD, theirs.Zone, "only its controller's cards come back")
	})
}
