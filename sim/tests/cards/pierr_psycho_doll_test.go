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
	pierrPsychoDollUID         = "e69867c3-41e2-4219-a3a9-29a0630c28e8"
	pierrPsychoDollAttackerUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	pierrPsychoDollBigUID      = "84e1b416-c2d5-4ae1-aca0-025651c6aa58" // Tri-horn Shepherd (5000)
)

func TestPierrPsychoDoll(t *testing.T) {
	t.Run("has the printed stats and static abilities", func(t *testing.T) {
		scn, _, defender := setupPierrPsychoDollTest(t)
		pierr := putPierrPsychoDollTestCardInBattlezone(t, scn, defender.Player, pierrPsychoDollUID)
		startPierrPsychoDollTestTurn(t, scn)

		assert.Equal(t, "Pierr, Psycho Doll", pierr.Name)
		assert.Equal(t, 1000, pierr.Power)
		assert.Equal(t, 2, pierr.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, pierr.Civs)
		assert.True(t, pierr.HasFamily(family.DeathPuppet))
		assert.True(t, pierr.HasCondition(cnd.Blocker))
		assert.True(t, pierr.HasCondition(cnd.Slayer))
		assert.True(t, pierr.HasCondition(cnd.CantAttackPlayers))
		assert.True(t, pierr.HasCondition(cnd.CantAttackCreatures))
	})

	t.Run("blocks an attack on its controller and trades through slayer", func(t *testing.T) {
		scn, attackerRef, defender := setupPierrPsychoDollTest(t)
		pierr := putPierrPsychoDollTestCardInBattlezone(t, scn, defender.Player, pierrPsychoDollUID)
		attacker := putPierrPsychoDollTestCardInBattlezone(t, scn, attackerRef.Player, pierrPsychoDollAttackerUID)
		startPierrPsychoDollTestTurn(t, scn)

		shields, err := defender.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)
		shieldCount := len(shields)

		action, err := scn.ActionAttackPlayer(attackerRef, attacker.ID)
		require.NoError(t, err)
		require.NotEmpty(t, action.Cards)
		require.NoError(t, scn.ResolveAttack(attackerRef, action.Cards[0].CardID))

		assert.Equal(t, match.GRAVEYARD, pierr.Zone, "Pierr loses the battle it was forced into")
		assert.Equal(t, match.GRAVEYARD, attacker.Zone, "slayer destroys the attacker after the battle")

		remainingShields, err := defender.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remainingShields, shieldCount, "a blocked attack breaks no shields")
	})

	t.Run("survives the forced block when it wins the battle", func(t *testing.T) {
		scn, attackerRef, defender := setupPierrPsychoDollTest(t)
		pierr := putPierrPsychoDollTestCardInBattlezone(t, scn, defender.Player, pierrPsychoDollUID)
		attacker := putPierrPsychoDollTestCardInBattlezone(t, scn, attackerRef.Player, pierrPsychoDollAttackerUID)
		startPierrPsychoDollTestTurn(t, scn)

		// A power boost large enough to win the battle proves the forced block is
		// a real battle and not an unconditional destruction of both creatures.
		pierr.AddUniqueSourceCondition(cnd.PowerAmplifier, 5000, "pierr_psycho_doll_test_setup")

		action, err := scn.ActionAttackPlayer(attackerRef, attacker.ID)
		require.NoError(t, err)
		require.NotEmpty(t, action.Cards)
		require.NoError(t, scn.ResolveAttack(attackerRef, action.Cards[0].CardID))

		assert.Equal(t, match.BATTLEZONE, pierr.Zone)
		assert.True(t, pierr.Tapped, "blocking taps the blocker")
		assert.Equal(t, match.GRAVEYARD, attacker.Zone)
	})

	t.Run("does not block while it is tapped", func(t *testing.T) {
		scn, attackerRef, defender := setupPierrPsychoDollTest(t)
		pierr := putPierrPsychoDollTestCardInBattlezone(t, scn, defender.Player, pierrPsychoDollUID)
		attacker := putPierrPsychoDollTestCardInBattlezone(t, scn, attackerRef.Player, pierrPsychoDollAttackerUID)
		startPierrPsychoDollTestTurn(t, scn)
		pierr.Tapped = true

		shields, err := defender.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		action, err := scn.ActionAttackPlayer(attackerRef, attacker.ID)
		require.NoError(t, err)
		require.NotEmpty(t, action.Cards)
		require.NoError(t, scn.ResolveAttack(attackerRef, action.Cards[0].CardID))

		assert.Equal(t, match.BATTLEZONE, pierr.Zone)
		assert.Equal(t, match.BATTLEZONE, attacker.Zone)
		remainingShields, err := defender.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remainingShields, shieldCount-1)
	})

	t.Run("cannot attack even when a tapped creature is available", func(t *testing.T) {
		scn, attackerRef, defender := setupPierrPsychoDollTest(t)
		// Pierr belongs to the player whose turn it will be so that it could
		// legally attack if it were not restricted.
		pierr := putPierrPsychoDollTestCardInBattlezone(t, scn, attackerRef.Player, pierrPsychoDollUID)
		target := putPierrPsychoDollTestCardInBattlezone(t, scn, defender.Player, pierrPsychoDollBigUID)
		startPierrPsychoDollTestTurn(t, scn)
		target.Tapped = true

		warningStart, err := scn.MessageCount(attackerRef)
		require.NoError(t, err)
		require.Error(t, scn.ActionAttackCreature(attackerRef, pierr.ID, target.ID), "the attack must never reach target selection")
		require.NoError(t, scn.WaitForMessage(attackerRef, warningStart, "warn"))

		assert.False(t, pierr.Tapped)
		assert.Equal(t, match.BATTLEZONE, target.Zone)
	})
}

func setupPierrPsychoDollTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	attacker := scn.Match.CurrentPlayer()
	defender := scn.Match.PlayerRef(scn.Match.Opponent(attacker.Player))

	return scn, attacker, defender
}

// startPierrPsychoDollTestTurn cycles a full turn so that the untap step has
// rebuilt every intrinsic condition and the attacker is free of summoning
// sickness, then hands the turn back to the original attacker.
func startPierrPsychoDollTestTurn(t *testing.T, scn *scenario.TestScenario) {
	t.Helper()

	attacker := scn.Match.CurrentPlayer()
	defender := scn.Match.PlayerRef(scn.Match.Opponent(attacker.Player))
	require.NoError(t, scn.ActionEndTurn(attacker))
	require.NoError(t, scn.ActionEndTurn(defender))
}

func putPierrPsychoDollTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, "pierr_psycho_doll_test_setup")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}
