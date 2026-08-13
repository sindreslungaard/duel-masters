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
	venomCapsuleUID      = "51ef593f-a840-4123-9a29-c2f11aa6dbd9"
	venomCapsuleSetupSrc = "venom_capsule_test_setup"
)

func TestVenomCapsule(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		capsule := putCardInBattlezone(t, scn, player.Player, venomCapsuleUID, venomCapsuleSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Venom Capsule", capsule.Name)
		assert.Equal(t, 1000, capsule.Power)
		assert.Equal(t, 2, capsule.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, capsule.Civs)
		assert.Equal(t, []string{civ.Darkness}, capsule.ManaRequirement)
		assert.True(t, capsule.HasFamily(family.BrainJacker))
		assert.True(t, capsule.HasCondition(cnd.SilentSkill))
	})

	t.Run("breaks one of the opponent's shields", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		capsule := putCardInBattlezone(t, scn, player.Player, venomCapsuleUID, venomCapsuleSetupSrc)
		capsule.Tapped = true

		passTurnToSelf(t, scn, player, opponent)

		// Counted here rather than before the turn handover, because the
		// opponent draws a card during the turn of their own in between.
		shieldsBefore, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shieldsBefore)
		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		useSilentSkill(t, scn, player)

		// Breaking a shield asks which one, face down.
		require.NoError(t, scn.SubmitAction(player, shieldsBefore[0].ID))
		settleTurn(t, scn)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		assert.Len(t, shields, len(shieldsBefore)-1)
		assert.Len(t, hand, len(handBefore)+1, "a broken shield goes to its owner's hand")
		assert.Equal(t, match.HAND, shieldsBefore[0].Zone)
		assert.True(t, capsule.Tapped)
	})

	t.Run("breaking is not attacking, so it cannot win the game", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		capsule := putCardInBattlezone(t, scn, player.Player, venomCapsuleUID, venomCapsuleSetupSrc)
		capsule.Tapped = true

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		for _, shield := range shields {
			_, err := opponent.Player.MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, venomCapsuleSetupSrc)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)
		useSilentSkill(t, scn, player)

		// No shields left to choose from, so nothing is asked and nothing
		// happens: an effect that breaks shields never wins on its own.
		assert.False(t, scn.Match.IsClosed())
		assert.True(t, capsule.Tapped)
	})

	t.Run("declining leaves the shields alone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		capsule := putCardInBattlezone(t, scn, player.Player, venomCapsuleUID, venomCapsuleSetupSrc)
		capsule.Tapped = true

		shieldsBefore, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		declineSilentSkill(t, scn, player)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore))
		assert.False(t, capsule.Tapped)
	})
}
