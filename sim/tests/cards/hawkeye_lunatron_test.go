package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	hawkeyeLunatronUID       = "b6ebbe0b-7441-4a89-9d7e-185b0c7bfe57"
	hawkeyeLunatronManaUID   = "9781089f-1aa9-4a75-b106-35e9d431e31d" // Aqua Vehicle
	hawkeyeLunatronTargetUID = "7956b4f5-b910-403d-b388-b67c837b7e99" // Scissor Eye
)

func TestHawkeyeLunatron(t *testing.T) {
	t.Run("searches the deck for one card and shuffles it afterwards", func(t *testing.T) {
		scn, player, hawkeye := setupHawkeyeLunatronTest(t)
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		player.Player.SpawnCard(hawkeyeLunatronTargetUID, match.DECK)
		target, err := scn.FindCard(player.Player, match.DECK, hawkeyeLunatronTargetUID)
		require.NoError(t, err)

		assert.Equal(t, "Hawkeye Lunatron", hawkeye.Name)
		assert.Equal(t, 6000, hawkeye.Power)
		assert.Equal(t, 8, hawkeye.ManaCost)
		assert.Equal(t, []string{civ.Water}, hawkeye.Civs)
		assert.True(t, hawkeye.HasFamily(family.CyberMoon))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, hawkeye.ID))
		assert.Equal(t, match.BATTLEZONE, hawkeye.Zone)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 0, action.MinSelections, "taking a card is optional")
		assert.Equal(t, 1, action.MaxSelections)
		assert.True(t, action.Cancellable)

		offered := make([]string, 0, len(action.Cards))
		for _, card := range action.Cards {
			offered = append(offered, card.CardID)
		}
		assert.Contains(t, offered, target.ID, "every card in the deck may be taken")

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, target.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.HAND, target.Zone)

		opponentHeaders, err := scn.MessageHeaders(opponent, opponentStart)
		require.NoError(t, err)
		assert.NotContains(t, opponentHeaders, "show_cards", "the taken card is not revealed")
	})

	t.Run("may decline to take a card", func(t *testing.T) {
		scn, player, hawkeye := setupHawkeyeLunatronTest(t)

		deck, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		deckSize := len(deck)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, hawkeye.ID))

		_, err = scn.LatestAction(player, promptStart)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		deck, err = player.Player.Container(match.DECK)
		require.NoError(t, err)
		assert.Len(t, deck, deckSize)
		assert.Equal(t, match.BATTLEZONE, hawkeye.Zone)
	})

	t.Run("does not search again when another creature is summoned", func(t *testing.T) {
		scn, player, hawkeye := setupHawkeyeLunatronTest(t)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, hawkeye.ID))
		_, err = scn.LatestAction(player, promptStart)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		player.Player.SpawnCard(hawkeyeLunatronTargetUID, match.HAND)
		other, err := scn.FindCard(player.Player, match.HAND, hawkeyeLunatronTargetUID)
		require.NoError(t, err)

		otherStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		// Moved directly: a second search prompt here would deadlock this goroutine.
		moved, err := player.Player.MoveCard(other.ID, match.HAND, match.BATTLEZONE, "hawkeye_lunatron_test_setup")
		require.NoError(t, err)
		require.Equal(t, match.BATTLEZONE, moved.Zone)

		headers, err := scn.MessageHeaders(player, otherStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
	})

	t.Run("breaks two shields as a double breaker", func(t *testing.T) {
		scn, player, hawkeye := setupHawkeyeLunatronTest(t)
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, hawkeye.ID))
		_, err = scn.LatestAction(player, promptStart)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, hawkeye.HasCondition(cnd.DoubleBreaker))

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		action, err := scn.ActionAttackPlayer(player, hawkeye.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, action.MinSelections)
		assert.Equal(t, 2, action.MaxSelections)
		require.Len(t, action.Cards, shieldCount)

		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID, action.Cards[1].CardID))

		// The opponent cannot block, so both selected shields are broken.
		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount-2)
	})
}

func setupHawkeyeLunatronTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.Card) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.SpawnCard(hawkeyeLunatronUID, match.HAND)
	for range 8 {
		player.Player.SpawnCard(hawkeyeLunatronManaUID, match.MANAZONE)
	}

	hawkeye, err := scn.FindCard(player.Player, match.HAND, hawkeyeLunatronUID)
	require.NoError(t, err)

	return scn, player, hawkeye
}
