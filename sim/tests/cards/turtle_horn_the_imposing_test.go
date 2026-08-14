package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	turtleHornUID = "1c6e6c4a-9108-4e88-9730-a9a12ddb2dbd"
	// Kamikaze, Chainsaw Warrior is a shield trigger creature, so playing it off
	// a broken shield opens no prompt of its own.
	turtleHornTriggerUID  = "0103d28a-1c07-4cc3-916e-1fd67ef9595a"
	turtleHornAttackerUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	turtleHornSrc         = "turtle_horn_the_imposing_test_setup"
)

func TestTurtleHornTheImposing(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, turtleHornUID, turtleHornSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Turtle Horn, the Imposing", 2000, 3, []string{civ.Nature})
		assert.True(t, card.HasFamily(family.HornedBeast))
	})

	t.Run("an opponent's shield trigger fetches a creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, turtleHornUID, turtleHornSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, turtleHornAttackerUID, turtleHornSrc)

		trigger, err := opponent.Player.SpawnCard(turtleHornTriggerUID, match.SHIELDZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(handBefore)

		triggerStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, trigger.ID))

		// The opponent is offered their broken shield trigger.
		_, err = scn.WaitForAction(opponent, triggerStart)
		require.NoError(t, err, "the shield trigger should have been offered")

		searchStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, trigger.ID))

		// Playing it is what Turtle Horn is watching for.
		search, err := scn.WaitForAction(player, searchStart)
		require.NoError(t, err, "the deck search should have been offered")
		require.NotEmpty(t, search.Cards)

		taken := search.Cards[0].CardID
		require.NoError(t, scn.SubmitAction(player, taken))
		settleTurn(t, scn)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount+1)
		assert.Equal(t, match.BATTLEZONE, trigger.Zone, "and the trigger itself was played")
	})

	t.Run("the search may take nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, turtleHornUID, turtleHornSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, turtleHornAttackerUID, turtleHornSrc)

		trigger, err := opponent.Player.SpawnCard(turtleHornTriggerUID, match.SHIELDZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(handBefore)

		triggerStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, trigger.ID))

		_, err = scn.WaitForAction(opponent, triggerStart)
		require.NoError(t, err)

		searchStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(opponent, trigger.ID))

		_, err = scn.WaitForAction(player, searchStart)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		settleTurn(t, scn)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount)
	})

	t.Run("a shield trigger the opponent declines fetches nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, turtleHornUID, turtleHornSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, turtleHornAttackerUID, turtleHornSrc)

		trigger, err := opponent.Player.SpawnCard(turtleHornTriggerUID, match.SHIELDZONE)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		triggerStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, trigger.ID))

		_, err = scn.WaitForAction(opponent, triggerStart)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(opponent))

		// Nothing was used, so nothing is asked of Turtle Horn's controller and
		// the loop settles on its own.
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, trigger.Zone, "an unused trigger stays in hand")
	})

	t.Run("a plain shield fetches nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, turtleHornUID, turtleHornSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, turtleHornAttackerUID, turtleHornSrc)

		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))
		settleTurn(t, scn)
	})
}
