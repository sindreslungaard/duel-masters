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
	squawkingLunatronUID = "c3380f40-8a9c-49f7-9c6b-e00ae2212481"
	squawkingSetupSrc    = "squawking_lunatron_test_setup"
)

func TestSquawkingLunatron(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, squawkingLunatronUID, squawkingSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, lunatron, "Squawking Lunatron", 4000, 5, []string{civ.Water})
		assert.True(t, lunatron.HasFamily(family.CyberMoon))
		assert.True(t, lunatron.HasCondition(cnd.SilentSkill))
	})

	t.Run("it returns up to three cards from the mana zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, squawkingLunatronUID, squawkingSetupSrc)
		lunatron.Tapped = true

		mana := make([]*match.Card, 0, 4)
		for range 4 {
			card, err := player.Player.SpawnCard(immortalBaronVorgUID, match.MANAZONE)
			require.NoError(t, err)
			mana = append(mana, card)
		}

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		answerInTurn(t, scn, player, mana[0].ID, mana[1].ID, mana[2].ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, mana[0].Zone)
		assert.Equal(t, match.HAND, mana[1].Zone)
		assert.Equal(t, match.HAND, mana[2].Zone)
		assert.Equal(t, match.MANAZONE, mana[3].Zone)
	})

	t.Run("up to three means the choice can be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, squawkingLunatronUID, squawkingSetupSrc)
		lunatron.Tapped = true

		card, err := player.Player.SpawnCard(immortalBaronVorgUID, match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		cancelInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.MANAZONE, card.Zone)
		assert.True(t, lunatron.Tapped)
	})
}
