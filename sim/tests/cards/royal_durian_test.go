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
	royalDurianUID   = "89304216-4cf9-4bb0-903f-55e7248655df"
	royalDurianSetup = "royal_durian_test_setup"
)

func TestRoyalDurian(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		durian := putCardInBattlezone(t, scn, player.Player, royalDurianUID, royalDurianSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, durian, "Royal Durian", 1000, 5, []string{civ.Nature})
		assert.True(t, durian.HasFamily(family.WildVeggies))
		assert.True(t, durian.HasCondition(cnd.SilentSkill))
	})

	t.Run("it puts a Dragon from the mana zone into the battle zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		durian := putCardInBattlezone(t, scn, player.Player, royalDurianUID, royalDurianSetup)
		durian.Tapped = true

		dragon, err := player.Player.SpawnCard(rollickingDragonID, match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.BATTLEZONE, dragon.Zone)
	})
}
