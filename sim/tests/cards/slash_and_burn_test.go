package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	slashAndBurnUID = "a1227d09-3baf-4cf3-ae11-2eebe659f1ff"
	slashSmallUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	slashBigUID     = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	slashAndBurnSrc = "slash_and_burn_test_setup"
)

func TestSlashAndBurn(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(slashAndBurnUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Slash and Burn", spell.Name)
		assert.Equal(t, 4, spell.ManaCost)
		assert.Equal(t, []string{civ.Darkness, civ.Fire}, spell.Civs)
	})

	t.Run("a destroyed creature costs its controller mana and a shield", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, player.Player, slashBigUID, slashAndBurnSrc)
		victim := putCardInBattlezone(t, scn, opponent.Player, slashSmallUID, slashAndBurnSrc)

		theirMana, err := opponent.Player.SpawnCard(slashSmallUID, match.MANAZONE)
		require.NoError(t, err)

		shieldsBefore, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		castSpell(t, scn, player, slashAndBurnUID)
		require.NoError(t, scn.WaitForEventLoop())

		// Killed in battle rather than by a direct Destroy call: the trigger
		// prompts the opponent, so the destruction has to be raised by the
		// event loop for the test to be able to answer it.
		victim.Tapped = true
		require.NoError(t, scn.ActionAttackCreature(player, attacker.ID, victim.ID))

		// One card in mana, so that choice resolves itself; the shield still
		// has to be picked.
		answerInTurn(t, scn, opponent, shieldsBefore[0].ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, victim.Zone)
		assert.Equal(t, match.GRAVEYARD, theirMana.Zone)
		assert.Equal(t, match.GRAVEYARD, shieldsBefore[0].Zone)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore)-1)
	})

	t.Run("the caster's own losses do not count", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		ownAttacker := putCardInBattlezone(t, scn, player.Player, slashSmallUID, slashAndBurnSrc)
		theirBig := putCardInBattlezone(t, scn, opponent.Player, slashBigUID, slashAndBurnSrc)

		theirMana, err := opponent.Player.SpawnCard(slashSmallUID, match.MANAZONE)
		require.NoError(t, err)

		castSpell(t, scn, player, slashAndBurnUID)
		require.NoError(t, scn.WaitForEventLoop())

		// Attacking into a bigger creature kills the attacker, which belongs to
		// the caster and so is none of the spell's business.
		theirBig.Tapped = true
		require.NoError(t, scn.ActionAttackCreature(player, ownAttacker.ID, theirBig.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, ownAttacker.Zone)
		assert.Equal(t, match.MANAZONE, theirMana.Zone, "the printed trigger is the opponent's creatures")
		assertShieldsIntact(t, scn, opponent)
	})

	t.Run("the effect is over on the following turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		attacker := putCardInBattlezone(t, scn, player.Player, slashBigUID, slashAndBurnSrc)
		victim := putCardInBattlezone(t, scn, opponent.Player, slashSmallUID, slashAndBurnSrc)

		theirMana, err := opponent.Player.SpawnCard(slashSmallUID, match.MANAZONE)
		require.NoError(t, err)

		castSpell(t, scn, player, slashAndBurnUID)
		require.NoError(t, scn.WaitForEventLoop())

		passTurnToSelf(t, scn, player, opponent)

		victim.Tapped = true
		require.NoError(t, scn.ActionAttackCreature(player, attacker.ID, victim.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, victim.Zone)
		assert.Equal(t, match.MANAZONE, theirMana.Zone, "the spell only lasted its own turn")
		assertShieldsIntact(t, scn, opponent)
	})
}

// assertShieldsIntact checks that a player still holds the five shields they
// started the duel with.
func assertShieldsIntact(t *testing.T, _ *scenario.TestScenario, player *match.PlayerReference) {
	t.Helper()

	shields, err := player.Player.Container(match.SHIELDZONE)
	require.NoError(t, err)
	assert.Len(t, shields, 5)
}
