package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	velyrikaDragonUID       = "0559379a-e0c5-4b59-bad6-c4bd07d73816"
	velyrikaDragonTargetUID = "91db2302-6794-4aa4-b17b-6637d356e9ac" // Astrocomet Dragon (Armored Dragon)
)

func TestVelyrikaDragon(t *testing.T) {
	t.Run("searches the deck for an Armored Dragon, shows it to the opponent, and shuffles afterwards", func(t *testing.T) {
		scn, player, velyrika := setupVelyrikaDragonTest(t)
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		player.Player.SpawnCard(velyrikaDragonTargetUID, match.DECK)
		target, err := scn.FindCard(player.Player, match.DECK, velyrikaDragonTargetUID)
		require.NoError(t, err)

		assert.Equal(t, "Velyrika Dragon", velyrika.Name)
		assert.Equal(t, 7000, velyrika.Power)
		assert.Equal(t, 7, velyrika.ManaCost)
		assert.Equal(t, []string{civ.Fire}, velyrika.Civs)
		assert.Equal(t, []string{civ.Fire}, velyrika.ManaRequirement)
		assert.True(t, velyrika.HasFamily(family.ArmoredDragon))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		opponentStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, velyrika.ID))
		assert.Equal(t, match.BATTLEZONE, velyrika.Zone)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 0, action.MinSelections, "taking an Armored Dragon is optional")
		assert.Equal(t, 1, action.MaxSelections)
		assert.True(t, action.Cancellable)

		offered := make([]string, 0, len(action.Cards))
		for _, card := range action.Cards {
			offered = append(offered, card.CardID)
		}
		assert.Contains(t, offered, target.ID, "every Armored Dragon in the deck may be taken")

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, target.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.HAND, target.Zone)

		opponentChat, err := scn.ChatMessages(opponent, opponentStart)
		require.NoError(t, err)
		found := false
		for _, message := range opponentChat {
			if strings.Contains(message, target.Name) {
				found = true
				break
			}
		}
		assert.True(t, found, "the taken Armored Dragon's name is announced to the opponent")
	})

	t.Run("may decline to take an Armored Dragon", func(t *testing.T) {
		scn, player, velyrika := setupVelyrikaDragonTest(t)

		player.Player.SpawnCard(velyrikaDragonTargetUID, match.DECK)

		deck, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		deckSize := len(deck)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, velyrika.ID))
		_, err = scn.LatestAction(player, promptStart)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		deck, err = player.Player.Container(match.DECK)
		require.NoError(t, err)
		assert.Len(t, deck, deckSize)
		assert.Equal(t, match.BATTLEZONE, velyrika.Zone)
	})

	t.Run("offers no selectable card when the deck has no Armored Dragon", func(t *testing.T) {
		scn, player, velyrika := setupVelyrikaDragonTest(t)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, velyrika.ID))

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 0, action.MinSelections)
		assert.Equal(t, 0, action.MaxSelections)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.BATTLEZONE, velyrika.Zone)
	})

	t.Run("does not search again when another creature is summoned", func(t *testing.T) {
		scn, player, velyrika := setupVelyrikaDragonTest(t)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, velyrika.ID))
		_, err = scn.LatestAction(player, promptStart)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		player.Player.SpawnCard(velyrikaDragonTargetUID, match.HAND)
		other, err := scn.FindCard(player.Player, match.HAND, velyrikaDragonTargetUID)
		require.NoError(t, err)

		otherStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		// Moved directly: a second search prompt here would deadlock this goroutine.
		moved, err := player.Player.MoveCard(other.ID, match.HAND, match.BATTLEZONE, "velyrika_dragon_test_setup")
		require.NoError(t, err)
		require.Equal(t, match.BATTLEZONE, moved.Zone)

		headers, err := scn.MessageHeaders(player, otherStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")
	})

	t.Run("breaks two shields as a double breaker", func(t *testing.T) {
		scn, player, velyrika := setupVelyrikaDragonTest(t)
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, velyrika.ID))
		_, err = scn.LatestAction(player, promptStart)
		require.NoError(t, err)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, velyrika.HasCondition(cnd.DoubleBreaker))

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		action, err := scn.ActionAttackPlayer(player, velyrika.ID)
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

func setupVelyrikaDragonTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.Card) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	player.Player.SpawnCard(velyrikaDragonUID, match.HAND)
	for range 7 {
		player.Player.SpawnCard(velyrikaDragonUID, match.MANAZONE)
	}

	velyrika, err := scn.FindCard(player.Player, match.HAND, velyrikaDragonUID)
	require.NoError(t, err)

	return scn, player, velyrika
}
