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
	rollickingTotemUID = "e278e681-ee2d-4865-be34-9933dca3d470"
	rollickingDragonID = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (Volcano Dragon)
	rollickingSetupSrc = "rollicking_totem_test_setup"
)

func TestRollickingTotem(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		totem := putCardInBattlezone(t, scn, player.Player, rollickingTotemUID, rollickingSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, totem, "Rollicking Totem", 4000, 5, []string{civ.Nature})
		assert.True(t, totem.HasFamily(family.MysteryTotem))
		assert.True(t, totem.HasCondition(cnd.SilentSkill))
	})

	t.Run("it puts a Dragon from the mana zone into the battle zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		totem := putCardInBattlezone(t, scn, player.Player, rollickingTotemUID, rollickingSetupSrc)
		totem.Tapped = true

		dragon, err := player.Player.SpawnCard(rollickingDragonID, match.MANAZONE)
		require.NoError(t, err)
		notADragon, err := player.Player.SpawnCard(immortalBaronVorgUID, match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		// The dragon is the only legal choice, so it is taken without asking.
		assert.Equal(t, match.BATTLEZONE, dragon.Zone)
		assert.True(t, dragon.HasCondition(cnd.SummoningSickness), "it was put into play, not summoned early")
		assert.Equal(t, match.MANAZONE, notADragon.Zone)
	})

	t.Run("a mana zone without dragons costs nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		totem := putCardInBattlezone(t, scn, player.Player, rollickingTotemUID, rollickingSetupSrc)
		totem.Tapped = true

		notADragon, err := player.Player.SpawnCard(immortalBaronVorgUID, match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.MANAZONE, notADragon.Zone)
		assert.True(t, totem.Tapped)
	})

	t.Run("it cannot reach into the opponent's mana zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		totem := putCardInBattlezone(t, scn, player.Player, rollickingTotemUID, rollickingSetupSrc)
		totem.Tapped = true

		theirDragon, err := opponent.Player.SpawnCard(rollickingDragonID, match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		assert.Equal(t, match.MANAZONE, theirDragon.Zone)
	})
}
