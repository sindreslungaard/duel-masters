package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const emergencyTyphoonUID = "ddf7ccd6-48e1-46f6-9800-367bf36ec29b"

func TestEmergencyTyphoon(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		spell, err := player.Player.SpawnCard(emergencyTyphoonUID, match.HAND)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Emergency Typhoon", spell.Name)
		assert.Equal(t, 2, spell.ManaCost)
		assert.Equal(t, []string{civ.Water}, spell.Civs)
		assert.True(t, spell.HasCondition(cnd.ShieldTrigger))
	})

	t.Run("drawing two and discarding one is a net gain of one card", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, emergencyTyphoonUID)

		answerDrawUpTo(t, scn, player, 2, true)
		answerInTurn(t, scn, player, handBefore[0].ID)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		// The spell left the hand, two cards were drawn and one discarded.
		assert.Len(t, hand, len(handBefore)+1)
		assert.Equal(t, match.GRAVEYARD, handBefore[0].Zone, "the chosen card is the one discarded")
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("the discard happens even when no cards are drawn", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		castSpell(t, scn, player, emergencyTyphoonUID)

		answerDrawUpTo(t, scn, player, 2, false)
		answerInTurn(t, scn, player, handBefore[0].ID)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		// "Draw up to 2 cards. Then discard a card" is not conditional on the
		// draw, so declining it entirely still costs a card.
		assert.Len(t, hand, len(handBefore)-1)
		assert.Equal(t, match.GRAVEYARD, handBefore[0].Zone)
	})

	t.Run("the spell cannot discard itself", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, emergencyTyphoonUID)

		answerDrawUpTo(t, scn, player, 2, false)

		// By the time the discard prompt opens, the spell must already have
		// left the hand for the graveyard as part of its own resolution, since
		// its "discard a card from your hand" text cannot refer to itself.
		assert.Equal(t, match.GRAVEYARD, spell.Zone)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err, "expected the discard prompt to be open")

		for _, offered := range action.Cards {
			assert.NotEqual(t, spell.ID, offered.CardID, "Emergency Typhoon must not be offered as its own discard target")
		}

		// Submitting the spell's own id anyway must be rejected like any other
		// out-of-list selection, not silently accepted as a second discard.
		invalidStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, spell.ID))
		require.NoError(t, scn.WaitForMessage(player, invalidStart, "action", "action_error", "warn", "wait", "state_update"))

		headers, err := scn.MessageHeaders(player, invalidStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action_error"), "an invalid selection should warn, not resolve")

		// The prompt is still open and answering it legitimately still works.
		answerInTurn(t, scn, player, handBefore[0].ID)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		assert.Len(t, hand, len(handBefore)-1)
		assert.Equal(t, match.GRAVEYARD, handBefore[0].Zone)
	})
}
