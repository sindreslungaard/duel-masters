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
	generalDarkFiendUID      = "fcf2f484-c471-4a5d-bc41-1fcd56604d73"
	generalDarkFiendSetupSrc = "general_dark_fiend_test_setup"
)

func TestGeneralDarkFiend(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		gdf := putCardInBattlezone(t, scn, player.Player, generalDarkFiendUID, generalDarkFiendSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, gdf, "General Dark Fiend", 6000, 5, []string{civ.Darkness})
		assert.True(t, gdf.HasFamily(family.DarkLord))
		assert.True(t, gdf.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("attacking makes its controller choose one of their own shields blind, and mills it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		gdf := putCardInBattlezone(t, scn, player.Player, generalDarkFiendUID, generalDarkFiendSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		myShields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		myShieldCount := len(myShields)
		myShieldIDs := make([]string, 0, myShieldCount)
		for _, s := range myShields {
			myShieldIDs = append(myShieldIDs, s.ID)
		}

		oppShields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		oppShieldCount := len(oppShields)

		millPrompt, err := scn.ActionAttackPlayer(player, gdf.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, millPrompt.MinSelections, "the choice is mandatory")
		assert.Equal(t, 1, millPrompt.MaxSelections)
		assert.False(t, millPrompt.Cancellable)

		offeredIDs := make([]string, 0, len(millPrompt.Cards))
		for _, c := range millPrompt.Cards {
			offeredIDs = append(offeredIDs, c.CardID)
		}
		assert.ElementsMatch(t, myShieldIDs, offeredIDs, "offered from its own controller's shields, not the opponent's")

		chosen := millPrompt.Cards[0].CardID

		shieldPromptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, chosen))

		shieldPrompt, err := scn.WaitForAction(player, shieldPromptStart)
		require.NoError(t, err, "expected the normal shield-selection prompt afterwards")

		require.NoError(t, scn.ResolveAttack(player, shieldPrompt.Cards[0].CardID, shieldPrompt.Cards[1].CardID))
		settleTurn(t, scn)

		milled, err := player.Player.GetCard(chosen, match.GRAVEYARD)
		require.NoError(t, err, "the chosen shield was moved to its controller's graveyard")
		assert.Equal(t, match.GRAVEYARD, milled.Zone)

		remainingMyShields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remainingMyShields, myShieldCount-1, "only the chosen shield left its controller's shield zone")

		remainingOppShields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remainingOppShields, oppShieldCount-2, "the attack itself still breaks 2 of the opponent's shields")
	})

	t.Run("with no shields of its own, nothing is asked and the attack proceeds normally", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		gdf := putCardInBattlezone(t, scn, player.Player, generalDarkFiendUID, generalDarkFiendSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		myShields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		for _, s := range myShields {
			_, err := player.Player.MoveCard(s.ID, match.SHIELDZONE, match.HAND, generalDarkFiendSetupSrc)
			require.NoError(t, err)
		}

		oppShields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		oppShieldCount := len(oppShields)

		shieldPrompt, err := scn.ActionAttackPlayer(player, gdf.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shieldPrompt.Cards[0].CardID, shieldPrompt.Cards[1].CardID))
		settleTurn(t, scn)

		remainingOppShields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remainingOppShields, oppShieldCount-2)
	})
}
