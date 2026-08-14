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
	tajimalVizierOfAquaUID = "1136d769-2ddf-4e56-ab54-6ed02e270452"
	// Magmadragon Melgars is a vanilla 4000 power fire creature. Without the
	// bonus that is an even trade which destroys both; with it Tajimal survives.
	tajimalFireUID              = "8112be9d-50a9-4489-b3f8-257aeed62205"
	tajimalNatureUID            = "84e1b416-c2d5-4ae1-aca0-025651c6aa58" // Tri-horn Shepherd (nature, 5000)
	tajimalVizierOfAquaSetupSrc = "tajimal_vizier_of_aqua_test_setup"
)

func TestTajimalVizierOfAqua(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupTajimalTest(t)
		tajimal := putCardInBattlezone(t, scn, player.Player, tajimalVizierOfAquaUID, tajimalVizierOfAquaSetupSrc)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.Equal(t, "Tajimal, Vizier of Aqua", tajimal.Name)
		assert.Equal(t, 4000, tajimal.Power)
		assert.Equal(t, 3, tajimal.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Water}, tajimal.Civs)
		assert.True(t, tajimal.HasFamily(family.Initiate))
		assert.True(t, tajimal.HasFamily(family.LiquidPeople))
		assert.True(t, tajimal.HasCondition(cnd.Blocker))
		assert.True(t, tajimal.HasCondition(cnd.CantAttackPlayers))
		assert.False(t, tajimal.HasCondition(cnd.CantAttackCreatures), "it may still attack creatures")
	})

	t.Run("survives a fire creature it would otherwise trade with", func(t *testing.T) {
		scn, defender, attacker := setupTajimalTest(t)
		tajimal := putCardInBattlezone(t, scn, defender.Player, tajimalVizierOfAquaUID, tajimalVizierOfAquaSetupSrc)
		fireAttacker := putCardInBattlezone(t, scn, attacker.Player, tajimalFireUID, tajimalVizierOfAquaSetupSrc)

		require.NoError(t, scn.ActionEndTurn(defender))
		tajimal.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(attacker, fireAttacker.ID, tajimal.ID))

		assert.Equal(t, match.GRAVEYARD, fireAttacker.Zone, "8000 beats 4000")
		assert.Equal(t, match.BATTLEZONE, tajimal.Zone, "an even trade without the bonus would destroy it too")
		assert.Equal(t, 4000, scn.Match.GetPower(tajimal, false), "the bonus applies only during the battle")
	})

	t.Run("gets no bonus against a creature of another civilization", func(t *testing.T) {
		scn, defender, attacker := setupTajimalTest(t)
		tajimal := putCardInBattlezone(t, scn, defender.Player, tajimalVizierOfAquaUID, tajimalVizierOfAquaSetupSrc)
		natureAttacker := putCardInBattlezone(t, scn, attacker.Player, tajimalNatureUID, tajimalVizierOfAquaSetupSrc)

		require.NoError(t, scn.ActionEndTurn(defender))
		tajimal.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(attacker, natureAttacker.ID, tajimal.ID))

		assert.Equal(t, match.GRAVEYARD, tajimal.Zone, "5000 beats an unboosted 4000")
		assert.Equal(t, match.BATTLEZONE, natureAttacker.Zone)
	})

	t.Run("gets the bonus while attacking a fire creature as well", func(t *testing.T) {
		scn, player, opponent := setupTajimalTest(t)
		tajimal := putCardInBattlezone(t, scn, player.Player, tajimalVizierOfAquaUID, tajimalVizierOfAquaSetupSrc)
		fireTarget := putCardInBattlezone(t, scn, opponent.Player, tajimalFireUID, tajimalVizierOfAquaSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))
		fireTarget.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(player, tajimal.ID, fireTarget.ID))

		assert.Equal(t, match.GRAVEYARD, fireTarget.Zone)
		assert.Equal(t, match.BATTLEZONE, tajimal.Zone)
	})

	t.Run("cannot attack the player", func(t *testing.T) {
		scn, player, opponent := setupTajimalTest(t)
		tajimal := putCardInBattlezone(t, scn, player.Player, tajimalVizierOfAquaUID, tajimalVizierOfAquaSetupSrc)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))

		_, err := scn.ActionAttackPlayer(player, tajimal.ID)
		assert.Error(t, err)
		assert.False(t, tajimal.Tapped, "a refused attack does not tap it")
	})
}

func setupTajimalTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	return scn, player, opponent
}
