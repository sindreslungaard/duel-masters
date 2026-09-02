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
	dynoMantisTheMightspinnerUID              = "7e9eddab-97ea-4b14-8e43-ede104dbf99a"
	dynoMantisTheMightspinnerAlly5000UID      = "c7fec5e8-4e56-451b-a7b6-ad08680703a4" // La Byle, Seeker of the Winds (power 5000)
	dynoMantisTheMightspinnerAllyWeakUID      = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (power 2000)
	dynoMantisTheMightspinnerQuixoticUID      = "701500bb-7ade-4d3d-9938-8367fa3a71bd" // Quixotic Hero Swine Snout (power 1000, +3000 per other creature summoned, until end of turn)
	dynoMantisTheMightspinnerPowerAttackerUID = "84595f9f-c9c1-4d50-8e8f-29e5ef63bfbf" // Sword Butterfly (power 2000, power attacker +3000)
	dynoMantisTheMightspinnerMuscleChargerUID = "1fbf62d7-aca7-4f0e-bf22-9170e36aad57" // Muscle Charger (+3000 to all own creatures until end of turn)
	dynoMantisTheMightspinnerSetupSrc         = "dyno_mantis_the_mightspinner_test_setup"
	dynoMantisTheMightspinnerThreshold        = 5000
)

func TestDynoMantisTheMightspinner(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dyno := putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerUID, dynoMantisTheMightspinnerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, dyno, "Dyno Mantis, the Mightspinner", 7000, 5, []string{civ.Nature})
		assert.True(t, dyno.HasFamily(family.GiantInsect))
		assert.True(t, dyno.HasCondition(cnd.Evolution))
		assert.True(t, dyno.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("it breaks two shields on its own when it attacks", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dyno := putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerUID, dynoMantisTheMightspinnerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		action, err := scn.ActionAttackPlayer(player, dyno.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, action.MaxSelections, "its own double breaker, not tripled by its own aura")

		require.NoError(t, scn.CancelAction(player))
	})

	t.Run("a controlled creature with power 5000 or more breaks one more shield", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerUID, dynoMantisTheMightspinnerSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerAlly5000UID, dynoMantisTheMightspinnerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		shieldsBefore, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(shieldsBefore), 2)

		action, err := scn.ActionAttackPlayer(player, ally.ID)
		require.NoError(t, err)
		require.Equal(t, 2, action.MaxSelections, "1 printed + 1 from Dyno Mantis (power is exactly 5000)")

		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID, action.Cards[1].CardID))
		settleTurn(t, scn)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore)-2)
	})

	t.Run("a controlled creature under 5000 power gets no bonus", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerUID, dynoMantisTheMightspinnerSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerAllyWeakUID, dynoMantisTheMightspinnerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		action, err := scn.ActionAttackPlayer(player, ally.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, action.MaxSelections)

		require.NoError(t, scn.CancelAction(player))
	})

	t.Run("it does not affect the opponent's creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerUID, dynoMantisTheMightspinnerSetupSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, dynoMantisTheMightspinnerAlly5000UID, dynoMantisTheMightspinnerSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		require.True(t, scn.Match.IsPlayerTurn(opponent.Player))

		action, err := scn.ActionAttackPlayer(opponent, theirs.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, action.MaxSelections, "Dyno Mantis only watches its controller's creatures")

		require.NoError(t, scn.CancelAction(opponent))
	})

	t.Run("the bonus leaves when Dyno Mantis leaves the battle zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dyno := putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerUID, dynoMantisTheMightspinnerSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerAlly5000UID, dynoMantisTheMightspinnerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		require.True(t, ally.HasCondition(cnd.ShieldBreakModifier), "sanity check: the bonus is active while Dyno Mantis is in play")

		_, err := player.Player.MoveCard(dyno.ID, match.BATTLEZONE, match.GRAVEYARD, dynoMantisTheMightspinnerSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.False(t, ally.HasCondition(cnd.ShieldBreakModifier))

		action, err := scn.ActionAttackPlayer(player, ally.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, action.MaxSelections, "the bonus left with Dyno Mantis")

		require.NoError(t, scn.CancelAction(player))
	})

	// The following three cases match published rulings for this card:
	//
	// Q: If Quixotic Hero Swine Snout gets 6000 power from 2 creatures being
	// put into the battle zone, does Dyno Mantis, the Mightspinner effect
	// also let it break an additional shield?
	// A: Yes. Quixotic's power is increased for the entire turn, not just
	// while attacking, and therefore it has 5000 power or more.
	t.Run("a temporary power boost from another creature entering counts for the whole turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerUID, dynoMantisTheMightspinnerSetupSrc)
		quixotic := putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerQuixoticUID, dynoMantisTheMightspinnerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		require.Less(t, scn.Match.GetPower(quixotic, false), dynoMantisTheMightspinnerThreshold, "starts under the threshold")

		// Two more creatures entering the battle zone push Quixotic to
		// 1000 + 3000 + 3000 = 7000, comfortably past the threshold.
		putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerAllyWeakUID, dynoMantisTheMightspinnerSetupSrc)
		putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerAllyWeakUID, dynoMantisTheMightspinnerSetupSrc)
		require.GreaterOrEqual(t, scn.Match.GetPower(quixotic, false), dynoMantisTheMightspinnerThreshold)

		action, err := scn.ActionAttackPlayer(player, quixotic.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, action.MaxSelections, "the boost lasts the whole turn, not just while attacking")

		require.NoError(t, scn.CancelAction(player))
	})

	// Q: Does Dyno Mantis, the Mightspinner effect affect creatures that
	// have "power attacker" that have 5000 or more power while attacking?
	// A: Yes. A creature with "power attacker" still has the increased power
	// until after the shields have been broken, so at the time of the
	// shield break, it still has 5000 power or more.
	t.Run("a power attacker bonus still counts at the moment shields are counted", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerUID, dynoMantisTheMightspinnerSetupSrc)
		attacker := putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerPowerAttackerUID, dynoMantisTheMightspinnerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		require.True(t, attacker.HasCondition(cnd.PowerAttacker))
		require.Less(t, scn.Match.GetPower(attacker, false), dynoMantisTheMightspinnerThreshold, "under the threshold while not attacking")
		require.GreaterOrEqual(t, scn.Match.GetPower(attacker, true), dynoMantisTheMightspinnerThreshold, "power attacker reaches the threshold while attacking")

		action, err := scn.ActionAttackPlayer(player, attacker.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, action.MaxSelections, "power attacker still applies at the point shields are counted")

		require.NoError(t, scn.CancelAction(player))
	})

	// Q: If Muscle Charger is cast while Dyno Mantis, the Mightspinner is in
	// the battle zone, do all creatures that have 5000 or more power get to
	// break an additional shield?
	// A: Yes. This is also true for other card effects that increase power
	// until the end of the turn.
	t.Run("a shared power boost from a spell counts too", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerUID, dynoMantisTheMightspinnerSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, dynoMantisTheMightspinnerAllyWeakUID, dynoMantisTheMightspinnerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		require.Less(t, scn.Match.GetPower(ally, false), dynoMantisTheMightspinnerThreshold)

		castSpell(t, scn, player, dynoMantisTheMightspinnerMuscleChargerUID)
		require.GreaterOrEqual(t, scn.Match.GetPower(ally, false), dynoMantisTheMightspinnerThreshold)

		action, err := scn.ActionAttackPlayer(player, ally.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, action.MaxSelections, "a shared, until-end-of-turn power boost still counts")

		require.NoError(t, scn.CancelAction(player))
	})
}
