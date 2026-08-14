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
	belixTheExplorerUID = "8ed1afbb-c05e-46c0-949b-7cc4b9bf7cee"
	belixSpellUID       = "5883180e-d88c-4f24-b17c-f5a837420147" // Terror Pit
	belixSetupSrc       = "belix_the_explorer_test_setup"
)

func TestBelixTheExplorer(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		belix := putCardInBattlezone(t, scn, player.Player, belixTheExplorerUID, belixSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, belix, "Belix, the Explorer", 3000, 2, []string{civ.Light})
		assert.True(t, belix.HasFamily(family.Gladiator))
		assert.True(t, belix.HasCondition(cnd.Blocker))
		assert.True(t, belix.HasCondition(cnd.CantAttackPlayers))
	})

	t.Run("summoning it returns a spell from the mana zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		spellInMana, err := player.Player.SpawnCard(belixSpellUID, match.MANAZONE)
		require.NoError(t, err)
		creatureInMana, err := player.Player.SpawnCard(immortalBaronVorgUID, match.MANAZONE)
		require.NoError(t, err)

		belix, err := player.Player.SpawnCard(belixTheExplorerUID, match.HAND)
		require.NoError(t, err)
		for range 2 {
			_, err := player.Player.SpawnCard(belixTheExplorerUID, match.MANAZONE)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionPlayCard(player, belix.ID))
		require.NoError(t, scn.WaitForEventLoop())

		// Only one spell is in the mana zone, so it comes back without asking.
		assert.Equal(t, match.HAND, spellInMana.Zone)
		assert.Equal(t, match.MANAZONE, creatureInMana.Zone, "creatures are not spells")
	})

	t.Run("a mana zone without spells costs nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		belix, err := player.Player.SpawnCard(belixTheExplorerUID, match.HAND)
		require.NoError(t, err)
		for range 2 {
			_, err := player.Player.SpawnCard(belixTheExplorerUID, match.MANAZONE)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.ActionPlayCard(player, belix.ID))
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment should ask anything")
		assert.Equal(t, match.BATTLEZONE, belix.Zone)
	})
}
