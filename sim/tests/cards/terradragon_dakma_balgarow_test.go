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
	terradragonDakmaBalgarowUID      = "8ae71daf-a39d-4bff-8ebc-b966bf6a059c"
	terradragonDakmaBalgarowFillerID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	terradragonDakmaBalgarowSetupSrc = "terradragon_dakma_balgarow_test_setup"
)

func TestTerradragonDakmaBalgarow(t *testing.T) {
	t.Run("the shield bonus is inactive outside the battle zone", func(t *testing.T) {
		scn, owner, _ := setupTerradragonDakmaBalgarowTest(t)
		owner.Player.SpawnCard(terradragonDakmaBalgarowUID, match.HAND)
		terradragon, err := scn.FindCard(owner.Player, match.HAND, terradragonDakmaBalgarowUID)
		require.NoError(t, err)

		assert.Equal(t, "Terradragon Dakma Balgarow", terradragon.Name)
		assert.Equal(t, 1000, terradragon.Power)
		assert.Equal(t, 7, terradragon.ManaCost)
		assert.Equal(t, []string{civ.Nature}, terradragon.Civs)
		assert.True(t, terradragon.HasFamily(family.EarthDragon))

		assert.Equal(t, 1000, scn.Match.GetPower(terradragon, false))
		assert.False(t, terradragon.HasCondition(cnd.DoubleBreaker))
		assert.False(t, terradragon.HasCondition(cnd.TripleBreaker))
	})

	t.Run("counts every shield both players have", func(t *testing.T) {
		scn, owner, opponent := setupTerradragonDakmaBalgarowTest(t)
		terradragon := putTerradragonDakmaBalgarowTestCardInBattlezone(t, scn, owner.Player)

		setTerradragonDakmaBalgarowTestShields(t, scn, owner, 1)
		setTerradragonDakmaBalgarowTestShields(t, scn, opponent, 1)
		assert.Equal(t, 5000, scn.Match.GetPower(terradragon, false))
		assert.Equal(t, 1000, terradragon.Power, "dynamic power must not mutate printed power")

		setTerradragonDakmaBalgarowTestShields(t, scn, opponent, 4)
		assert.Equal(t, 11000, scn.Match.GetPower(terradragon, false))

		setTerradragonDakmaBalgarowTestShields(t, scn, owner, 4)
		assert.Equal(t, 17000, scn.Match.GetPower(terradragon, false))
	})

	t.Run("switches between no breaker, double breaker and triple breaker", func(t *testing.T) {
		scn, owner, opponent := setupTerradragonDakmaBalgarowTest(t)
		terradragon := putTerradragonDakmaBalgarowTestCardInBattlezone(t, scn, owner.Player)

		// 2 shields -> 5000 power
		setTerradragonDakmaBalgarowTestShields(t, scn, owner, 1)
		setTerradragonDakmaBalgarowTestShields(t, scn, opponent, 1)
		assert.False(t, terradragon.HasCondition(cnd.DoubleBreaker))
		assert.False(t, terradragon.HasCondition(cnd.TripleBreaker))

		// 3 shields -> 7000 power
		setTerradragonDakmaBalgarowTestShields(t, scn, owner, 2)
		assert.True(t, terradragon.HasCondition(cnd.DoubleBreaker))
		assert.False(t, terradragon.HasCondition(cnd.TripleBreaker))

		// 6 shields -> 13000 power, still short of the triple breaker threshold
		setTerradragonDakmaBalgarowTestShields(t, scn, opponent, 4)
		assert.True(t, terradragon.HasCondition(cnd.DoubleBreaker))
		assert.False(t, terradragon.HasCondition(cnd.TripleBreaker))

		// 7 shields -> exactly 15000 power
		setTerradragonDakmaBalgarowTestShields(t, scn, owner, 3)
		assert.True(t, terradragon.HasCondition(cnd.TripleBreaker))
		assert.False(t, terradragon.HasCondition(cnd.DoubleBreaker), "triple breaker replaces double breaker")

		// Back below 6000 power
		setTerradragonDakmaBalgarowTestShields(t, scn, owner, 1)
		setTerradragonDakmaBalgarowTestShields(t, scn, opponent, 1)
		assert.False(t, terradragon.HasCondition(cnd.DoubleBreaker))
		assert.False(t, terradragon.HasCondition(cnd.TripleBreaker))
	})

	t.Run("the breakers survive a turn boundary", func(t *testing.T) {
		scn, owner, opponent := setupTerradragonDakmaBalgarowTest(t)
		terradragon := putTerradragonDakmaBalgarowTestCardInBattlezone(t, scn, owner.Player)

		setTerradragonDakmaBalgarowTestShields(t, scn, owner, 2)
		setTerradragonDakmaBalgarowTestShields(t, scn, opponent, 1)
		require.True(t, terradragon.HasCondition(cnd.DoubleBreaker))

		require.NoError(t, scn.ActionEndTurn(owner))
		assert.True(t, terradragon.HasCondition(cnd.DoubleBreaker), "still a double breaker on the opponent's turn")
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.True(t, terradragon.HasCondition(cnd.DoubleBreaker))
		assert.False(t, terradragon.HasCondition(cnd.TripleBreaker))
	})

	t.Run("breaks three shields at 15000 power and drops back to two", func(t *testing.T) {
		scn, owner, opponent := setupTerradragonDakmaBalgarowTest(t)
		terradragon := putTerradragonDakmaBalgarowTestCardInBattlezone(t, scn, owner.Player)

		// 2 own + 5 opposing shields -> 15000 power.
		setTerradragonDakmaBalgarowTestShields(t, scn, owner, 2)
		setTerradragonDakmaBalgarowTestShields(t, scn, opponent, 5)
		require.Equal(t, 15000, scn.Match.GetPower(terradragon, false))
		require.True(t, terradragon.HasCondition(cnd.TripleBreaker))

		action, err := scn.ActionAttackPlayer(owner, terradragon.ID)
		require.NoError(t, err)
		assert.Equal(t, 3, action.MinSelections)
		assert.Equal(t, 3, action.MaxSelections)

		require.NoError(t, scn.ResolveAttack(owner, action.Cards[0].CardID, action.Cards[1].CardID, action.Cards[2].CardID))

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.Len(t, remaining, 2)

		// 2 own + 2 opposing shields -> 9000 power.
		assert.Equal(t, 9000, scn.Match.GetPower(terradragon, false))
		assert.True(t, terradragon.HasCondition(cnd.DoubleBreaker))
		assert.False(t, terradragon.HasCondition(cnd.TripleBreaker))
	})

	t.Run("loses its breakers when it leaves the battle zone", func(t *testing.T) {
		scn, owner, _ := setupTerradragonDakmaBalgarowTest(t)
		terradragon := putTerradragonDakmaBalgarowTestCardInBattlezone(t, scn, owner.Player)
		require.True(t, terradragon.HasCondition(cnd.TripleBreaker), "10 shields grant 21000 power")

		moved, err := owner.Player.MoveCard(terradragon.ID, match.BATTLEZONE, match.GRAVEYARD, terradragonDakmaBalgarowSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		assert.False(t, terradragon.HasCondition(cnd.DoubleBreaker))
		assert.False(t, terradragon.HasCondition(cnd.TripleBreaker))
		assert.Equal(t, 1000, scn.Match.GetPower(terradragon, false))
	})
}

func setupTerradragonDakmaBalgarowTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	owner := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))

	return scn, owner, opponent
}

func putTerradragonDakmaBalgarowTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player) *match.Card {
	t.Helper()

	player.SpawnCard(terradragonDakmaBalgarowUID, match.HAND)
	card, err := scn.FindCard(player, match.HAND, terradragonDakmaBalgarowUID)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, terradragonDakmaBalgarowSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}

// setTerradragonDakmaBalgarowTestShields brings a player's shield zone to an
// exact size, moving shields through the engine so the continuous effect is
// re-evaluated for every change.
func setTerradragonDakmaBalgarowTestShields(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, target int) {
	t.Helper()

	for {
		shields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		if len(shields) == target {
			return
		}

		if len(shields) > target {
			moved, err := player.Player.MoveCard(shields[len(shields)-1].ID, match.SHIELDZONE, match.GRAVEYARD, terradragonDakmaBalgarowSetupSrc)
			require.NoError(t, err)
			require.Equal(t, match.GRAVEYARD, moved.Zone)
			continue
		}

		player.Player.SpawnCard(terradragonDakmaBalgarowFillerID, match.HAND)
		filler, err := scn.FindCard(player.Player, match.HAND, terradragonDakmaBalgarowFillerID)
		require.NoError(t, err)
		moved, err := player.Player.MoveCard(filler.ID, match.HAND, match.SHIELDZONE, terradragonDakmaBalgarowSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.SHIELDZONE, moved.Zone)
	}
}
