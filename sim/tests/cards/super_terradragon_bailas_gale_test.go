package cards

import (
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	superTerradragonBailasGaleUID = "44cbf9fd-3906-4032-939c-f702ffda7415"
	// Ghost Touch (darkness, shield trigger): the opponent discards a random
	// card, with no selection of its own, so it resolves without opening a
	// second prompt.
	bailasGaleShieldTriggerSpellUID = "ae797f95-54b1-48e9-9216-f315b39826bd"
	bailasGaleSetupSrc              = "super_terradragon_bailas_gale_test_setup"
)

func TestSuperTerradragonBailasGale(t *testing.T) {
	t.Run("returns a spell cast by its own shield trigger to hand instead of the graveyard", func(t *testing.T) {
		scn, defender, attacker := setupDuel(t)

		bailasGale := putCardInBattlezone(t, scn, defender.Player, superTerradragonBailasGaleUID, bailasGaleSetupSrc)
		shieldSpell, err := defender.Player.SpawnCard(bailasGaleShieldTriggerSpellUID, match.SHIELDZONE)
		require.NoError(t, err)
		raider := putCardInBattlezone(t, scn, attacker.Player, scowlingTomatoUID, bailasGaleSetupSrc)

		// A single untap step, on the way into attacker's turn, gives Bailas
		// Gale, the shield spell, and the raider their conditions, and clears
		// the raider's summoning sickness.
		require.NoError(t, scn.ActionEndTurn(defender))

		_, err = scn.ActionAttackPlayer(attacker, raider.ID)
		require.NoError(t, err)

		promptStart, err := scn.MessageCount(defender)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(attacker, shieldSpell.ID))

		action, err := scn.WaitForAction(defender, promptStart)
		require.NoError(t, err, "expected the shield trigger prompt offering the spell")
		offeredCardIDs := make([]string, 0, len(action.Cards))
		for _, offered := range action.Cards {
			offeredCardIDs = append(offeredCardIDs, offered.CardID)
		}
		require.Contains(t, offeredCardIDs, shieldSpell.ID)

		completionStart, err := scn.MessageCount(defender)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(defender, shieldSpell.ID))
		require.NoError(t, scn.WaitForMessage(defender, completionStart, "state_update"))

		assert.Equal(t, match.HAND, shieldSpell.Zone, "Bailas Gale should return the spell to hand instead of the graveyard")
		assert.Equal(t, match.BATTLEZONE, bailasGale.Zone)
	})

	t.Run("does not affect a spell cast normally from hand", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		putCardInBattlezone(t, scn, player.Player, superTerradragonBailasGaleUID, bailasGaleSetupSrc)

		spell := castSpell(t, scn, player, bailasGaleShieldTriggerSpellUID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, spell.Zone, "a normal cast is unaffected by the shield-trigger-only effect")
	})
}
