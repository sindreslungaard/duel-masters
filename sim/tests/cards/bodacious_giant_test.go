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
	bodaciousGiantUID          = "7a39340a-4601-4c34-8754-05554d49cbf4"
	bodaciousGiantAttackerUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	bodaciousGiantDecoyUID     = "abe034cb-5b1b-4a9f-9519-0af2bfcc4796" // Cragsaur
	bodaciousGiantCantAttackID = "a15ddc75-f015-42b6-be15-eb17e1da2779" // Battery Cluster (can attack nothing)
	bodaciousGiantSetupSrc     = "bodacious_giant_test_setup"
)

func TestBodaciousGiant(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, _, _, giant, _ := setupTauntingBodaciousGiant(t)

		assert.Equal(t, "Bodacious Giant", giant.Name)
		assert.Equal(t, 12000, giant.Power)
		assert.Equal(t, 8, giant.ManaCost)
		assert.Equal(t, []string{civ.Nature}, giant.Civs)
		assert.True(t, giant.HasFamily(family.Giant))
		assert.True(t, giant.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("forces the opponent's attacks onto itself while tapped", func(t *testing.T) {
		scn, defender, attacker, giant, attackers := setupTauntingBodaciousGiant(t, bodaciousGiantAttackerUID, bodaciousGiantDecoyUID)
		creature := attackers[0]

		// A second tapped creature on the defending side is not a legal target.
		decoy := attackers[1]
		moved, err := attacker.Player.MoveCard(decoy.ID, match.BATTLEZONE, match.GRAVEYARD, bodaciousGiantSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
		defenderDecoy := putBodaciousGiantTestCardInBattlezone(t, scn, defender.Player, bodaciousGiantDecoyUID)
		defenderDecoy.Tapped = true

		warningStart, err := scn.MessageCount(attacker)
		require.NoError(t, err)
		_, err = scn.ActionAttackPlayer(attacker, creature.ID)
		require.Error(t, err, "attacking the player must be rejected")

		warnings, err := scn.Warnings(attacker, warningStart)
		require.NoError(t, err)
		require.NotEmpty(t, warnings)
		assert.Contains(t, warnings[0], "Bodacious Giant", "the warning names the creature that must be attacked")
		assert.False(t, creature.Tapped, "a refused attack does not tap the attacker")

		shields, err := defender.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		require.Error(t, scn.ActionAttackCreature(attacker, creature.ID, defenderDecoy.ID), "another tapped creature is not offered")
		require.NoError(t, scn.ActionAttackCreature(attacker, creature.ID, giant.ID))

		assert.Equal(t, match.GRAVEYARD, creature.Zone, "2000 power loses to 12000")
		assert.Equal(t, match.BATTLEZONE, giant.Zone)
		assert.Equal(t, match.BATTLEZONE, defenderDecoy.Zone)
		remaining, err := defender.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount, "no shields are broken")
	})

	t.Run("stops taunting once it has been attacked that turn", func(t *testing.T) {
		scn, defender, attacker, giant, attackers := setupTauntingBodaciousGiant(t, bodaciousGiantAttackerUID, bodaciousGiantAttackerUID)

		require.NoError(t, scn.ActionAttackCreature(attacker, attackers[0].ID, giant.ID))
		require.Equal(t, match.GRAVEYARD, attackers[0].Zone)

		shields, err := defender.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		action, err := scn.ActionAttackPlayer(attacker, attackers[1].ID)
		require.NoError(t, err, "attacking the player is allowed once the giant has been attacked")
		require.NoError(t, scn.ResolveAttack(attacker, action.Cards[0].CardID))

		remaining, err := defender.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount-1)
	})

	t.Run("the taunt returns on the next turn", func(t *testing.T) {
		scn, defender, attacker, giant, attackers := setupTauntingBodaciousGiant(t, bodaciousGiantAttackerUID, bodaciousGiantAttackerUID)

		require.NoError(t, scn.ActionAttackCreature(attacker, attackers[0].ID, giant.ID))
		require.Equal(t, match.GRAVEYARD, attackers[0].Zone)

		require.NoError(t, scn.ActionEndTurn(attacker))
		require.NoError(t, scn.ActionEndTurn(defender))
		giant.Tapped = true

		_, err := scn.ActionAttackPlayer(attacker, attackers[1].ID)
		assert.Error(t, err, "a new turn resets the attacked flag")
	})

	t.Run("does not taunt while it is untapped", func(t *testing.T) {
		scn, defender, attacker, giant, attackers := setupTauntingBodaciousGiant(t, bodaciousGiantAttackerUID)
		giant.Tapped = false

		shields, err := defender.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shields)

		action, err := scn.ActionAttackPlayer(attacker, attackers[0].ID)
		require.NoError(t, err)
		require.NoError(t, scn.ResolveAttack(attacker, action.Cards[0].CardID))

		remaining, err := defender.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, shieldCount-1)
	})

	t.Run("does not taunt on its own controller's turn", func(t *testing.T) {
		scn, defender, attacker, giant, attackers := setupTauntingBodaciousGiant(t, bodaciousGiantAttackerUID)
		own := putBodaciousGiantTestCardInBattlezone(t, scn, defender.Player, bodaciousGiantAttackerUID)

		// Clear the taunt so the attacker may end their turn.
		require.NoError(t, scn.ActionAttackCreature(attacker, attackers[0].ID, giant.ID))
		require.NoError(t, scn.ActionEndTurn(attacker))
		giant.Tapped = true

		action, err := scn.ActionAttackPlayer(defender, own.ID)
		require.NoError(t, err, "a tapped Bodacious Giant never restricts its own controller")
		require.NoError(t, scn.ResolveAttack(defender, action.Cards[0].CardID))

		remaining, err := attacker.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, 4)
	})

	t.Run("holds up the turn while an able creature has not attacked", func(t *testing.T) {
		scn, _, attacker, _, _ := setupTauntingBodaciousGiant(t, bodaciousGiantAttackerUID)

		warningStart, err := scn.MessageCount(attacker)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(attacker))
		assert.True(t, scn.Match.IsPlayerTurn(attacker.Player), "the turn cannot end yet")

		warnings, err := scn.Warnings(attacker, warningStart)
		require.NoError(t, err)
		mustAttack := 0
		for _, warning := range warnings {
			if strings.Contains(warning, "must attack before you can end your turn") {
				mustAttack++
			}
		}
		assert.Equal(t, 1, mustAttack)
	})

	t.Run("ignores creatures that could never attack it", func(t *testing.T) {
		scn, _, attacker, _, attackers := setupTauntingBodaciousGiant(t, bodaciousGiantAttackerUID, bodaciousGiantCantAttackID)
		attackers[0].Tapped = true

		// Battery Cluster can attack neither players nor creatures, so it is never
		// able to attack the giant and must not hold up the turn.
		require.True(t, attackers[1].HasCondition(cnd.CantAttackCreatures))

		require.NoError(t, scn.ActionEndTurn(attacker))
		assert.False(t, scn.Match.IsPlayerTurn(attacker.Player))
	})
}

// setupTauntingBodaciousGiant puts a Bodacious Giant into the battle zone of the
// player who moves first and the given creatures into their opponent's, cycles
// to the opponent's turn so every intrinsic condition has been rebuilt, then
// taps the giant so its taunt is active.
func setupTauntingBodaciousGiant(t *testing.T, attackerUIDs ...string) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference, *match.Card, []*match.Card) {
	t.Helper()

	scn := scenario.New()
	defender := scn.Match.CurrentPlayer()
	attacker := scn.Match.PlayerRef(scn.Match.Opponent(defender.Player))

	giant := putBodaciousGiantTestCardInBattlezone(t, scn, defender.Player, bodaciousGiantUID)

	attackers := make([]*match.Card, 0, len(attackerUIDs))
	for _, uid := range attackerUIDs {
		attackers = append(attackers, putBodaciousGiantTestCardInBattlezone(t, scn, attacker.Player, uid))
	}

	require.NoError(t, scn.ActionEndTurn(defender))
	giant.Tapped = true

	return scn, defender, attacker, giant, attackers
}

func putBodaciousGiantTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, bodaciousGiantSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}
