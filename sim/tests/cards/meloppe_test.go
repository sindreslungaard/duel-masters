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
	meloppeUID      = "0451a36e-fe88-4817-ad94-3dbd9e460fb8"
	meloppeSetupSrc = "meloppe_test_setup"

	// meloppeVanillaAttackerUID is Lah, Purification Enforcer: a plain 5500
	// power creature with no effects of its own, used as an attacker that
	// isn't Meloppe itself.
	meloppeVanillaAttackerUID = "0cc5279e-0a26-41a8-a2a5-f7711120b772"

	// meloppeDoubleBreakerUID is Astrocomet Dragon: double breaker with no
	// other effect, used to confirm the shield count survives the swap.
	meloppeDoubleBreakerUID = "91db2302-6794-4aa4-b17b-6637d356e9ac"
)

func TestMeloppe(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, meloppeUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Meloppe", 1000, 3, []string{civ.Water})
		assert.True(t, card.HasFamily(family.CyberLord))
	})

	t.Run("defender chooses their own shields instead of the attacker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// The defender controls Meloppe; the attacker is an unrelated vanilla creature.
		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, meloppeVanillaAttackerUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		attackerStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionAttackPlayerAsync(player, attacker.ID))

		action, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err, "expected the defender, not the attacker, to be prompted for shields")
		assert.False(t, action.Cancellable, "a chooser Meloppe installed can't call off the attacker's attack")
		assert.Equal(t, 1, action.MinSelections)
		assert.Equal(t, 1, action.MaxSelections)
		assert.Len(t, action.Cards, len(shields))

		// The attacker never received a shield-selection prompt of their own,
		// only the "waiting for your opponent" notice.
		attackerHeaders, err := scn.MessageHeaders(player, attackerStart)
		require.NoError(t, err)
		assert.NotContains(t, attackerHeaders, "action")

		chosen := shields[0]
		require.NoError(t, scn.SubmitAction(opponent, chosen.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, chosen.Zone)
	})

	t.Run("attacking with Meloppe itself still leaves the choice to its opponent", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		meloppe := putCardInBattlezone(t, scn, player.Player, meloppeUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionAttackPlayerAsync(player, meloppe.ID))

		action, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err, "expected the defender to be prompted even though they don't control Meloppe")
		assert.False(t, action.Cancellable)

		chosen := shields[0]
		require.NoError(t, scn.SubmitAction(opponent, chosen.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, chosen.Zone)
	})

	t.Run("double breaker still asks the defender for two shields", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, meloppeDoubleBreakerUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(shields), 2)
		require.True(t, attacker.HasCondition(cnd.DoubleBreaker))

		defenderStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		require.NoError(t, scn.ActionAttackPlayerAsync(player, attacker.ID))

		action, err := scn.WaitForAction(opponent, defenderStart)
		require.NoError(t, err)
		assert.Equal(t, 2, action.MinSelections)
		assert.Equal(t, 2, action.MaxSelections)

		require.NoError(t, scn.SubmitAction(opponent, shields[0].ID, shields[1].ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, shields[0].Zone)
		assert.Equal(t, match.HAND, shields[1].Zone)
	})

	t.Run("Meloppe leaving play beforehand restores the attacker as chooser", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		meloppe := putCardInBattlezone(t, scn, opponent.Player, meloppeUID, meloppeSetupSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, meloppeVanillaAttackerUID, meloppeSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		_, err := opponent.Player.MoveCard(meloppe.ID, match.BATTLEZONE, match.GRAVEYARD, meloppeSetupSrc)
		require.NoError(t, err)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		attackerStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.ActionAttackPlayerAsync(player, attacker.ID))

		action, err := scn.WaitForAction(player, attackerStart)
		require.NoError(t, err, "expected the attacker to be prompted again once Meloppe left play")
		assert.True(t, action.Cancellable)

		require.NoError(t, scn.SubmitAction(player, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, shields[0].Zone)
	})
}
