package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Immortal Baron, Vorg is a vanilla 2000 creature and the whole test deck, so
	// attacks resolve without blockers or shield triggers interfering.
	attackStepAttackerUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5"
	// Popple, Flowerpetal Dancer's tap ability puts the top card of its
	// controller's deck into their manazone, which resolves without a prompt of
	// its own. It does not count as charging mana.
	attackStepTapAbilityUID = "bb3a50dd-049d-488d-b89c-779dcf29b82e"
	attackStepSetupSrc      = "attack_step_test_setup"
)

// The attack step begins when an attack is confirmed, not when one is declared.
// Whether a player may still summon, cast or charge mana hangs off that
// distinction, and every step on the way to a confirmed attack is cancellable:
// shield selection, attack target selection, and choosing between several tap
// abilities. Backing out of any of them has to leave the turn exactly as it was.
//
// Shield selection is the case worth guarding hardest, because it is cancelled
// on a nested context of its own. Cancellation does not propagate to a parent
// context, so the outer AttackPlayer context still looks perfectly healthy after
// the player has walked away, and anything that reads it to decide whether the
// attack happened gets the answer wrong.
func TestAttackStepBeginsOnlyOnConfirmedAttacks(t *testing.T) {

	t.Run("cancelling shield selection leaves the turn untouched", func(t *testing.T) {
		scn, player, opponent := attackStepScenario(t)
		attacker := putInBattleZone(t, scn, player, attackStepAttackerUID)
		passTurn(t, scn, player, opponent)

		_, err := scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		assert.False(t, attacker.Tapped, "an abandoned attack does not tap the attacker")
		assert.False(t, inAttackStep(scn), "the match must not have entered the attack step")
		assert.True(t, player.Player.CanChargeMana)

		// The flag is only worth anything if the action it guards still works.
		spare := firstHandCard(t, player)
		require.NoError(t, scn.ActionChargeMana(player, spare.ID))
		assert.Equal(t, match.MANAZONE, spare.Zone, "mana can still be charged after backing out")
	})

	t.Run("cancelling shield selection leaves the creature able to attack again", func(t *testing.T) {
		scn, player, opponent := attackStepScenario(t)
		attacker := putInBattleZone(t, scn, player, attackStepAttackerUID)
		passTurn(t, scn, player, opponent)

		_, err := scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		_, err = scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err, "the same creature should be able to declare a second attack")
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())

		shieldsAfter, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shieldsAfter, shieldCount-1, "the second attack should break a shield")
		assert.True(t, inAttackStep(scn))
	})

	t.Run("confirming an attack on the player enters the attack step", func(t *testing.T) {
		scn, player, opponent := attackStepScenario(t)
		attacker := putInBattleZone(t, scn, player, attackStepAttackerUID)
		passTurn(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, attacker.Tapped, "a confirmed attack taps the attacker")
		assert.True(t, inAttackStep(scn))
		assert.False(t, player.Player.CanChargeMana)

		spare := firstHandCard(t, player)
		start, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.ActionChargeMana(player, spare.ID))
		assert.Equal(t, match.HAND, spare.Zone, "charging mana after an attack is refused")

		warnings, err := scn.Warnings(player, start)
		require.NoError(t, err)
		assert.True(t, containsWarning(warnings, "can't charge mana"), "the refusal should be explained, got %v", warnings)
	})

	t.Run("cancelling attack target selection leaves the turn untouched", func(t *testing.T) {
		scn, player, opponent := attackStepScenario(t)
		attacker := putInBattleZone(t, scn, player, attackStepAttackerUID)
		defender := putInBattleZone(t, scn, opponent, attackStepAttackerUID)
		passTurn(t, scn, player, opponent)

		// Only a tapped creature can be attacked, and this one never acted.
		defender.Tapped = true

		_, err := scn.ActionAttackCreaturePrompt(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.CancelAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		assert.False(t, attacker.Tapped)
		assert.Equal(t, match.BATTLEZONE, defender.Zone, "no battle should have happened")
		assert.False(t, inAttackStep(scn))
		assert.True(t, player.Player.CanChargeMana)

		spare := firstHandCard(t, player)
		require.NoError(t, scn.ActionChargeMana(player, spare.ID))
		assert.Equal(t, match.MANAZONE, spare.Zone)
	})

	t.Run("confirming an attack on a creature enters the attack step", func(t *testing.T) {
		scn, player, opponent := attackStepScenario(t)
		attacker := putInBattleZone(t, scn, player, attackStepAttackerUID)
		defender := putInBattleZone(t, scn, opponent, attackStepAttackerUID)
		passTurn(t, scn, player, opponent)

		defender.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(player, attacker.ID, defender.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, inAttackStep(scn))
		assert.False(t, player.Player.CanChargeMana)
	})

	t.Run("using a tap ability enters the attack step", func(t *testing.T) {
		scn, player, opponent := attackStepScenario(t)
		popple := putInBattleZone(t, scn, player, attackStepTapAbilityUID)

		// A creature conjured straight into play has no conditions until an untap
		// step runs, so the tap ability itself only exists after passing a turn.
		passTurn(t, scn, player, opponent)
		require.True(t, popple.HasCondition(cnd.TapAbility))

		require.NoError(t, scn.ActionUseTapAbility(player, popple.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, popple.Tapped)
		assert.True(t, inAttackStep(scn), "tap abilities are only usable during the attack step")
		assert.False(t, player.Player.CanChargeMana)

		spare := firstHandCard(t, player)
		require.NoError(t, scn.ActionChargeMana(player, spare.ID))
		assert.Equal(t, match.HAND, spare.Zone, "charging mana after a tap ability is refused")
	})

	t.Run("a rejected attack never enters the attack step", func(t *testing.T) {
		scn, player, opponent := attackStepScenario(t)
		attacker := putInBattleZone(t, scn, player, attackStepAttackerUID)
		passTurn(t, scn, player, opponent)

		attacker.Tapped = true

		_, err := scn.ActionAttackPlayer(player, attacker.ID)
		require.Error(t, err, "a tapped creature cannot attack")

		assert.False(t, inAttackStep(scn))
		assert.True(t, player.Player.CanChargeMana)
	})

	t.Run("the flags are cleared again on the next turn", func(t *testing.T) {
		scn, player, opponent := attackStepScenario(t)
		attacker := putInBattleZone(t, scn, player, attackStepAttackerUID)
		passTurn(t, scn, player, opponent)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		_, err = scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))
		require.NoError(t, scn.WaitForEventLoop())
		require.True(t, inAttackStep(scn))

		passTurn(t, scn, player, opponent)

		assert.False(t, inAttackStep(scn), "a new turn starts outside the attack step")
		assert.True(t, player.Player.CanChargeMana)
		assert.False(t, player.Player.HasChargedMana)

		spare := firstHandCard(t, player)
		require.NoError(t, scn.ActionChargeMana(player, spare.ID))
		assert.Equal(t, match.MANAZONE, spare.Zone)
	})
}

func attackStepScenario(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}

// inAttackStep reads the step the same way Match.BroadcastState does, so the
// assertions do not depend on which non-attack step the match happens to be in.
func inAttackStep(scn *scenario.TestScenario) bool {
	_, ok := scn.Match.Step.(*match.AttackStep)
	return ok
}

func firstHandCard(t *testing.T, player *match.PlayerReference) *match.Card {
	t.Helper()

	hand, err := player.Player.Container(match.HAND)
	require.NoError(t, err)
	require.NotEmpty(t, hand, "the player needs a card in hand to charge mana with")

	return hand[0]
}

func containsWarning(warnings []string, substring string) bool {
	for _, warning := range warnings {
		if strings.Contains(warning, substring) {
			return true
		}
	}

	return false
}
