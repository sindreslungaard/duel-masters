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
	chargeWhipperUID        = "9700f13c-8491-4543-8137-0f1316bd86ac"
	chargeWhipperTriggerUID = "5883180e-d88c-4f24-b17c-f5a837420147" // Terror Pit (shield trigger)
	chargeWhipperPlainUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	chargeWhipperSetupSrc   = "charge_whipper_test_setup"
)

func TestChargeWhipper(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		whipper := putCardInBattlezone(t, scn, player.Player, chargeWhipperUID, chargeWhipperSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Charge Whipper", whipper.Name)
		assert.Equal(t, 2000, whipper.Power)
		assert.Equal(t, 3, whipper.ManaCost)
		assert.Equal(t, []string{civ.Water}, whipper.Civs)
		assert.Equal(t, []string{civ.Water}, whipper.ManaRequirement)
		assert.True(t, whipper.HasFamily(family.CyberVirus))
		assert.True(t, whipper.HasCondition(cnd.SilentSkill))
	})

	t.Run("swaps a card from hand for one of the shields", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		whipper := putCardInBattlezone(t, scn, player.Player, chargeWhipperUID, chargeWhipperSetupSrc)
		whipper.Tapped = true

		fromHand, err := player.Player.SpawnCard(chargeWhipperPlainUID, match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shieldsBefore)
		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		takenShield := shieldsBefore[0]

		useSilentSkill(t, scn, player)
		require.NoError(t, scn.SubmitAction(player, fromHand.ID))
		require.NoError(t, scn.SubmitAction(player, takenShield.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.SHIELDZONE, fromHand.Zone)
		assert.Equal(t, match.HAND, takenShield.Zone)

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore), "one in, one out")
		assert.Len(t, hand, len(handBefore)+1, "the swap is even, so only the draw step changes the count")
		assert.True(t, whipper.Tapped)
	})

	t.Run("the shield taken back never offers its shield trigger", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		whipper := putCardInBattlezone(t, scn, player.Player, chargeWhipperUID, chargeWhipperSetupSrc)
		whipper.Tapped = true

		fromHand, err := player.Player.SpawnCard(chargeWhipperPlainUID, match.HAND)
		require.NoError(t, err)

		// Replace one shield with a card that would be very noticeable if it
		// were cast: breaking it would destroy a creature.
		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)
		_, err = player.Player.MoveCard(shields[0].ID, match.SHIELDZONE, match.GRAVEYARD, chargeWhipperSetupSrc)
		require.NoError(t, err)
		trigger, err := player.Player.SpawnCard(chargeWhipperTriggerUID, match.SHIELDZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		useSilentSkill(t, scn, player)
		require.NoError(t, scn.SubmitAction(player, fromHand.ID))
		require.NoError(t, scn.SubmitAction(player, trigger.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, trigger.Zone)

		// Two prompts belong to the card itself, the hand pick and the shield
		// pick. A cast shield trigger would have opened a third.
		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.LessOrEqual(t, countHeaders(headers, "action"), 2, "the shield is moved, not broken")
	})

	t.Run("declining the first half skips the second", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		whipper := putCardInBattlezone(t, scn, player.Player, chargeWhipperUID, chargeWhipperSetupSrc)
		whipper.Tapped = true

		_, err := player.Player.SpawnCard(chargeWhipperPlainUID, match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		useSilentSkill(t, scn, player)
		require.NoError(t, scn.CancelAction(player))
		settleTurn(t, scn)

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore), "no shield comes back when none went in")
		assert.True(t, whipper.Tapped)
	})

	t.Run("an empty hand means nothing is swapped", func(t *testing.T) {
		scn, player, opponent := setupSilentSkillTest(t)
		whipper := putCardInBattlezone(t, scn, player.Player, chargeWhipperUID, chargeWhipperSetupSrc)
		whipper.Tapped = true

		passTurnToSelf(t, scn, player, opponent)

		// Emptied while the event loop is parked on the silent skill prompt, so
		// the ability resolves against a hand with nothing to offer.
		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		for _, card := range hand {
			_, err := player.Player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, chargeWhipperSetupSrc)
			require.NoError(t, err)
		}

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		useSilentSkill(t, scn, player)

		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore))
		assert.True(t, whipper.Tapped)
	})
}

func countHeaders(headers []string, wanted string) int {
	count := 0
	for _, header := range headers {
		if header == wanted {
			count++
		}
	}

	return count
}
