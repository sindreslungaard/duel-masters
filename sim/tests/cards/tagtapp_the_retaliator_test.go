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
	tagtappTheRetaliatorUID       = "6c1ef042-f5f6-45b9-849f-a3c1eafbec1e"
	tagtappTheRetaliatorMelniaUID = "ddccdc18-92ef-431e-913e-71ba5bb6b1b1" // Melnia, the Aqua Shadow (water/darkness)
	tagtappTheRetaliatorSetupSrc  = "tagtapp_the_retaliator_test_setup"
)

func TestTagtappTheRetaliator(t *testing.T) {
	t.Run("the bonus is inactive outside the battle zone", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		tagtapp := spawnMulticolorCardInHand(t, scn, player, tagtappTheRetaliatorUID)
		spawnCivMana(t, opponent, civ.Water, 4)

		assert.Equal(t, "Tagtapp, the Retaliator", tagtapp.Name)
		assert.Equal(t, 3000, tagtapp.Power)
		assert.Equal(t, 3, tagtapp.ManaCost)
		assert.Equal(t, []string{civ.Fire, civ.Nature}, tagtapp.Civs)
		assert.True(t, tagtapp.HasFamily(family.SpiritQuartz))

		assert.Equal(t, 3000, scn.Match.GetPower(tagtapp, false))
		assert.False(t, tagtapp.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("gains power for each water card in the opponent's mana zone", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		tagtapp := putCardInBattlezone(t, scn, player.Player, tagtappTheRetaliatorUID, tagtappTheRetaliatorSetupSrc)
		assert.Equal(t, 3000, scn.Match.GetPower(tagtapp, false))

		spawnCivMana(t, opponent, civ.Water, 2)
		assert.Equal(t, 5000, scn.Match.GetPower(tagtapp, false))

		spawnCivMana(t, opponent, civ.Fire, 3)
		assert.Equal(t, 5000, scn.Match.GetPower(tagtapp, false), "only water cards count")

		spawnCivMana(t, player, civ.Water, 3)
		assert.Equal(t, 5000, scn.Match.GetPower(tagtapp, false), "only the opponent's mana zone counts")
	})

	t.Run("a multicolored water card in the opponent's mana zone counts", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		tagtapp := putCardInBattlezone(t, scn, player.Player, tagtappTheRetaliatorUID, tagtappTheRetaliatorSetupSrc)

		melnia, err := opponent.Player.SpawnCard(tagtappTheRetaliatorMelniaUID, match.MANAZONE)
		require.NoError(t, err)
		require.True(t, melnia.HasCiv(civ.Water))
		require.True(t, melnia.HasCiv(civ.Darkness))

		assert.Equal(t, 4000, scn.Match.GetPower(tagtapp, false), "water/darkness is a water card")
	})

	t.Run("becomes a double breaker at 6000 power and loses it again", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		tagtapp := putCardInBattlezone(t, scn, player.Player, tagtappTheRetaliatorUID, tagtappTheRetaliatorSetupSrc)

		water := spawnCivMana(t, opponent, civ.Water, 2)
		// Spawning does not run through the engine, so force a state refresh.
		require.NoError(t, scn.ActionEndTurn(player))
		assert.Equal(t, 5000, scn.Match.GetPower(tagtapp, false))
		assert.False(t, tagtapp.HasCondition(cnd.DoubleBreaker))

		spawnCivMana(t, opponent, civ.Water, 1)
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.Equal(t, 6000, scn.Match.GetPower(tagtapp, false))
		assert.True(t, tagtapp.HasCondition(cnd.DoubleBreaker))
		assert.False(t, tagtapp.HasCondition(cnd.TripleBreaker), "it has no triple breaker tier")

		moved, err := opponent.Player.MoveCard(water[0].ID, match.MANAZONE, match.GRAVEYARD, tagtappTheRetaliatorSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
		assert.Equal(t, 5000, scn.Match.GetPower(tagtapp, false))
		assert.False(t, tagtapp.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("breaks two shields while it has 6000 power", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

		tagtapp := putCardInBattlezone(t, scn, player.Player, tagtappTheRetaliatorUID, tagtappTheRetaliatorSetupSrc)
		spawnCivMana(t, opponent, civ.Water, 3)
		require.NoError(t, scn.ActionEndTurn(player))
		require.NoError(t, scn.ActionEndTurn(opponent))
		require.Equal(t, 6000, scn.Match.GetPower(tagtapp, false))

		action, err := scn.ActionAttackPlayer(player, tagtapp.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, action.MinSelections)
		assert.Equal(t, 2, action.MaxSelections)

		require.NoError(t, scn.ResolveAttack(player, action.Cards[0].CardID, action.Cards[1].CardID))

		remaining, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remaining, 3)
	})
}
