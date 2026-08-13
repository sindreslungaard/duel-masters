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
	burnwispLizardUID           = "affde51a-afb7-49e7-9f85-b9cfee945523"
	burnwispLizardSilentUID     = "460fc2eb-c7cd-42d5-9bed-a98de4f59026" // Milporo (silent skill)
	burnwispLizardOrdinaryUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	burnwispLizardManaUID       = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (water mana)
	burnwispLizardTestSetupName = "burnwisp_lizard_test_setup"
)

func TestBurnwispLizard(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lizard := putCardInBattlezone(t, scn, player.Player, burnwispLizardUID, burnwispLizardTestSetupName)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Burnwisp Lizard", lizard.Name)
		assert.Equal(t, 4000, lizard.Power)
		assert.Equal(t, 5, lizard.ManaCost)
		assert.Equal(t, []string{civ.Fire}, lizard.Civs)
		assert.Equal(t, []string{civ.Fire}, lizard.ManaRequirement)
		assert.True(t, lizard.HasFamily(family.MeltWarrior))
	})

	t.Run("gives speed attacker to its controller's silent skill creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, burnwispLizardUID, burnwispLizardTestSetupName)

		silent := putCardInBattlezone(t, scn, player.Player, burnwispLizardSilentUID, burnwispLizardTestSetupName)
		ordinary := putCardInBattlezone(t, scn, player.Player, burnwispLizardOrdinaryUID, burnwispLizardTestSetupName)
		theirSilent := putCardInBattlezone(t, scn, opponent.Player, burnwispLizardSilentUID, burnwispLizardTestSetupName)

		passTurnToSelf(t, scn, player, opponent)

		assert.True(t, silent.HasCondition(cnd.SpeedAttacker))
		assert.False(t, ordinary.HasCondition(cnd.SpeedAttacker), "only silent skill creatures gain it")
		assert.False(t, theirSilent.HasCondition(cnd.SpeedAttacker), "only its controller's creatures gain it")
	})

	t.Run("a silent skill creature summoned this turn can attack", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, burnwispLizardUID, burnwispLizardTestSetupName)

		for range 4 {
			_, err := player.Player.SpawnCard(burnwispLizardManaUID, match.MANAZONE)
			require.NoError(t, err)
		}

		_, err := player.Player.SpawnCard(burnwispLizardSilentUID, match.HAND)
		require.NoError(t, err)
		silent, err := scn.FindCard(player.Player, match.HAND, burnwispLizardSilentUID)
		require.NoError(t, err)

		// An untap step has to run before the card counts as having silent
		// skill, which is how a real card reaches hand: it was in the deck for
		// every untap step of the game. A card conjured straight into hand by
		// SpawnCard has no conditions at all until one runs.
		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionPlayCard(player, silent.ID))
		require.NoError(t, scn.WaitForEventLoop())
		require.Equal(t, match.BATTLEZONE, silent.Zone)

		require.True(t, silent.HasCondition(cnd.SummoningSickness))
		require.True(t, silent.HasCondition(cnd.SpeedAttacker))

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.NotEmpty(t, shields)

		_, err = scn.ActionAttackPlayer(player, silent.ID)
		require.NoError(t, err, "speed attacker should let it attack the turn it arrived")
		require.NoError(t, scn.ResolveAttack(player, shields[0].ID))

		assert.True(t, silent.Tapped, "attacking taps it, which sets up its silent skill next turn")
	})

	t.Run("the grant is removed when it leaves the battle zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lizard := putCardInBattlezone(t, scn, player.Player, burnwispLizardUID, burnwispLizardTestSetupName)
		silent := putCardInBattlezone(t, scn, player.Player, burnwispLizardSilentUID, burnwispLizardTestSetupName)

		passTurnToSelf(t, scn, player, opponent)
		require.True(t, silent.HasCondition(cnd.SpeedAttacker))

		_, err := player.Player.MoveCard(lizard.ID, match.BATTLEZONE, match.GRAVEYARD, burnwispLizardTestSetupName)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.False(t, silent.HasCondition(cnd.SpeedAttacker), "the bonus goes with its source")
	})

	t.Run("two copies each grant their own, and removing one keeps the other", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		first := putCardInBattlezone(t, scn, player.Player, burnwispLizardUID, burnwispLizardTestSetupName)
		putCardInBattlezone(t, scn, player.Player, burnwispLizardUID, burnwispLizardTestSetupName)
		silent := putCardInBattlezone(t, scn, player.Player, burnwispLizardSilentUID, burnwispLizardTestSetupName)

		passTurnToSelf(t, scn, player, opponent)
		require.True(t, silent.HasCondition(cnd.SpeedAttacker))

		_, err := player.Player.MoveCard(first.ID, match.BATTLEZONE, match.GRAVEYARD, burnwispLizardTestSetupName)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, silent.HasCondition(cnd.SpeedAttacker), "the surviving copy still grants it")
	})

	t.Run("a silent skill creature that arrives later is covered too", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, burnwispLizardUID, burnwispLizardTestSetupName)

		waiting, err := player.Player.SpawnCard(burnwispLizardSilentUID, match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)
		require.True(t, waiting.HasCondition(cnd.SilentSkill))

		silent, err := player.Player.MoveCard(waiting.ID, match.HAND, match.BATTLEZONE, burnwispLizardTestSetupName)
		require.NoError(t, err)
		require.Equal(t, match.BATTLEZONE, silent.Zone)
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, silent.HasCondition(cnd.SpeedAttacker), "the grant follows the battle zone, not a fixed list")
	})
}
