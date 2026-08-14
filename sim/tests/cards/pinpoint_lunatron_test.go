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
	pinpointLunatronUID      = "df56ab26-86cb-43cc-8d80-9afefc2b5162"
	pinpointLunatronPlainUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	pinpointLunatronSetupSrc = "pinpoint_lunatron_test_setup"
)

func TestPinpointLunatron(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, pinpointLunatronUID, pinpointLunatronSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Pinpoint Lunatron", lunatron.Name)
		assert.Equal(t, 2000, lunatron.Power)
		assert.Equal(t, 6, lunatron.ManaCost)
		assert.Equal(t, []string{civ.Water}, lunatron.Civs)
		assert.Equal(t, []string{civ.Water}, lunatron.ManaRequirement)
		assert.True(t, lunatron.HasFamily(family.CyberMoon))
		assert.True(t, lunatron.HasCondition(cnd.SilentSkill))
	})

	t.Run("returns an opposing creature to its owner's hand", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, pinpointLunatronUID, pinpointLunatronSetupSrc)
		lunatron.Tapped = true

		theirs := putCardInBattlezone(t, scn, opponent.Player, pinpointLunatronPlainUID, pinpointLunatronSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		require.NoError(t, scn.SubmitAction(player, theirs.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, theirs.Zone)
		assert.Equal(t, opponent.Player, theirs.Player, "it goes to its owner's hand, not the controller's")
		assert.True(t, lunatron.Tapped)
	})

	t.Run("returns a card from the opponent's mana zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, pinpointLunatronUID, pinpointLunatronSetupSrc)
		lunatron.Tapped = true

		theirMana, err := opponent.Player.SpawnCard(pinpointLunatronPlainUID, match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		require.NoError(t, scn.SubmitAction(player, theirMana.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, theirMana.Zone)
	})

	t.Run("it can return a card from its controller's own mana zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, pinpointLunatronUID, pinpointLunatronSetupSrc)
		lunatron.Tapped = true

		ownMana, err := player.Player.SpawnCard(pinpointLunatronPlainUID, match.MANAZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		require.NoError(t, scn.SubmitAction(player, ownMana.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, ownMana.Zone, "either player's mana zone is a legal source")
	})

	t.Run("it can return itself", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, pinpointLunatronUID, pinpointLunatronSetupSrc)
		lunatron.Tapped = true

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		// Its own controller's battle zone is eligible and it is the only
		// creature there, so the mandatory choice lands on itself.
		assert.Equal(t, match.HAND, lunatron.Zone)
	})

	t.Run("the choice is mandatory", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, pinpointLunatronUID, pinpointLunatronSetupSrc)
		lunatron.Tapped = true

		theirs := putCardInBattlezone(t, scn, opponent.Player, pinpointLunatronPlainUID, pinpointLunatronSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		// Cancelling is refused, so the prompt stays open until a real
		// selection arrives.
		cancelInTurn(t, scn, player)
		answerInTurn(t, scn, player, theirs.ID)
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, theirs.Zone)
	})

	t.Run("declining the silent skill returns nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lunatron := putCardInBattlezone(t, scn, player.Player, pinpointLunatronUID, pinpointLunatronSetupSrc)
		lunatron.Tapped = true

		theirs := putCardInBattlezone(t, scn, opponent.Player, pinpointLunatronPlainUID, pinpointLunatronSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		declineSilentSkill(t, scn, player)

		assert.Equal(t, match.BATTLEZONE, theirs.Zone)
		assert.Equal(t, match.BATTLEZONE, lunatron.Zone)
		assert.False(t, lunatron.Tapped)
	})
}
