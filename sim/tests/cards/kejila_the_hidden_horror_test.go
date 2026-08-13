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
	kejilaTheHiddenHorrorUID      = "1a2317fa-3710-418e-9b99-a3489f34285a"
	kejilaTheHiddenHorrorSetupSrc = "kejila_the_hidden_horror_test_setup"
)

func TestKejilaTheHiddenHorror(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		kejila := putCardInBattlezone(t, scn, player.Player, kejilaTheHiddenHorrorUID, kejilaTheHiddenHorrorSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Kejila, the Hidden Horror", kejila.Name)
		assert.Equal(t, 6000, kejila.Power)
		assert.Equal(t, 6, kejila.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, kejila.Civs)
		assert.Equal(t, []string{civ.Darkness}, kejila.ManaRequirement)
		assert.True(t, kejila.HasFamily(family.PandorasBox))
		assert.True(t, kejila.HasCondition(cnd.DoubleBreaker))
		assert.True(t, kejila.HasCondition(cnd.SilentSkill))
	})

	t.Run("breaks two of the opponent's shields", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		kejila := putCardInBattlezone(t, scn, player.Player, kejilaTheHiddenHorrorUID, kejilaTheHiddenHorrorSetupSrc)
		kejila.Tapped = true

		passTurnToSelf(t, scn, player, opponent)

		shieldsBefore, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(shieldsBefore), 2)
		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		useSilentSkill(t, scn, player)
		require.NoError(t, scn.SubmitAction(player, shieldsBefore[0].ID, shieldsBefore[1].ID))
		settleTurn(t, scn)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		assert.Len(t, shields, len(shieldsBefore)-2)
		assert.Len(t, hand, len(handBefore)+2)
		assert.Equal(t, match.HAND, shieldsBefore[0].Zone)
		assert.Equal(t, match.HAND, shieldsBefore[1].Zone)
		assert.True(t, kejila.Tapped)
	})

	t.Run("a single remaining shield is broken on its own", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		kejila := putCardInBattlezone(t, scn, player.Player, kejilaTheHiddenHorrorUID, kejilaTheHiddenHorrorSetupSrc)
		kejila.Tapped = true

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		for _, shield := range shields[1:] {
			_, err := opponent.Player.MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, kejilaTheHiddenHorrorSetupSrc)
			require.NoError(t, err)
		}

		lastShield := shields[0]

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		// Asking for two shields when only one exists offers the one that is
		// there rather than refusing to break anything.
		require.NoError(t, scn.SubmitAction(player, lastShield.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, lastShield.Zone)

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Empty(t, remaining)
		assert.False(t, scn.Match.IsClosed(), "emptying the shields this way does not win the game")
	})

	t.Run("declining leaves the shields alone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		kejila := putCardInBattlezone(t, scn, player.Player, kejilaTheHiddenHorrorUID, kejilaTheHiddenHorrorSetupSrc)
		kejila.Tapped = true

		passTurnToSelf(t, scn, player, opponent)

		shieldsBefore, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		declineSilentSkill(t, scn, player)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore))
		assert.False(t, kejila.Tapped)
	})
}
