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
	baleskBajUID        = "4dbc114f-b53e-45f5-baa5-5ca33a5a3337"
	baleskBajBlockerUID = "c7fec5e8-4e56-451b-a7b6-ad08680703a4" // La Byle, Seeker of the Winds (Blocker)
	baleskBajSetupSrc   = "balesk_baj_test_setup"
)

func TestBaleskBajTheTimeburner(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, baleskBajUID, baleskBajSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Balesk Baj, the Timeburner", 8000, 9, []string{civ.Fire})
		assert.True(t, card.HasFamily(family.ArmoredWyvern))
		assert.True(t, card.HasCondition(cnd.Evolution))
		assert.True(t, card.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("an unblocked attack returns it to hand and grants an extra turn", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, player.Player, baleskBajUID, baleskBajSetupSrc)

		action, err := scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID))

		require.NoError(t, scn.ActionEndTurn(player))
		assert.Equal(t, match.HAND, attacker.Zone)
		assert.True(t, scn.Match.IsPlayerTurn(player.Player), "the extra turn keeps the same player on turn")
	})

	t.Run("a blocked attack does not grant an extra turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, player.Player, baleskBajUID, baleskBajSetupSrc)
		blocker := putCardInBattlezone(t, scn, opponent.Player, baleskBajBlockerUID, baleskBajSetupSrc)
		// Grant the Blocker keyword directly: a real untap step for the
		// opponent would require ending player's turn first, which would
		// bounce Balesk Baj back to hand before it ever got to attack.
		blocker.AddUniqueSourceCondition(cnd.Blocker, true, blocker.ID)

		action, err := scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID))
		require.NoError(t, scn.SubmitAction(opponent, blocker.ID))
		require.NoError(t, scn.WaitForEventLoop())

		require.NoError(t, scn.ActionEndTurn(player))
		assert.Equal(t, match.HAND, attacker.Zone, "the return-to-hand clause is unconditional")
		assert.False(t, scn.Match.IsPlayerTurn(player.Player), "a blocked attack must not grant an extra turn")
	})

	t.Run("being destroyed after an unblocked attack does not cancel the extra turn", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, player.Player, baleskBajUID, baleskBajSetupSrc)

		action, err := scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID))
		require.Equal(t, match.BATTLEZONE, attacker.Zone)

		// Simulate being destroyed by a shield trigger or other effect after
		// the unblocked attack, but before the turn ends.
		moved, err := player.Player.MoveCard(attacker.ID, match.BATTLEZONE, match.GRAVEYARD, baleskBajSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
		require.NoError(t, scn.WaitForEventLoop())

		require.NoError(t, scn.ActionEndTurn(player))
		assert.True(t, scn.Match.IsPlayerTurn(player.Player), "the extra turn already triggered at the unblocked attack")
	})

	t.Run("not attacking grants no extra turn", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, player.Player, baleskBajUID, baleskBajSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		assert.Equal(t, match.HAND, attacker.Zone)
		assert.False(t, scn.Match.IsPlayerTurn(player.Player))
	})

	t.Run("the granted extra turn does not itself grant another one", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, player.Player, baleskBajUID, baleskBajSetupSrc)

		action, err := scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID))

		require.NoError(t, scn.ActionEndTurn(player))
		require.True(t, scn.Match.IsPlayerTurn(player.Player), "the extra turn was granted")

		// On the granted extra turn itself, Balesk Baj already left the
		// battle zone and does not attack again, so ending it must hand the
		// turn over normally rather than granting a second extra turn.
		require.NoError(t, scn.ActionEndTurn(player))
		assert.False(t, scn.Match.IsPlayerTurn(player.Player), "the flag must not leak into the extra turn itself")
	})
}
