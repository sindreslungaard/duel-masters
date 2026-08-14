package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	nialVizierOfDexterityUID = "4247022f-f630-4e8b-9e9f-9afaf5c61107"
	nialSetupSrc             = "nial_vizier_of_dexterity_test_setup"
)

func TestNialVizierOfDexterity(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		nial := putCardInBattlezone(t, scn, player.Player, nialVizierOfDexterityUID, nialSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, nial, "Nial, Vizier of Dexterity", 2500, 3, []string{civ.Light})
		assert.True(t, nial.HasFamily(family.Initiate))
	})

	t.Run("it may untap itself at the end of its controller's turn", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		nial := putCardInBattlezone(t, scn, player.Player, nialVizierOfDexterityUID, nialSetupSrc)
		nial.Tapped = true

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(player))

		_, err = scn.WaitForAction(player, promptStart)
		require.NoError(t, err)

		answerInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		assert.False(t, nial.Tapped)
	})

	t.Run("declining leaves it tapped", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		nial := putCardInBattlezone(t, scn, player.Player, nialVizierOfDexterityUID, nialSetupSrc)
		nial.Tapped = true

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(player))

		_, err = scn.WaitForAction(player, promptStart)
		require.NoError(t, err)

		cancelInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, nial.Tapped)
	})

	t.Run("an untapped creature is not offered anything", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, nialVizierOfDexterityUID, nialSetupSrc)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(player))

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
	})
}
