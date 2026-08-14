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

const terradragonCusdalfUID = "45b557c2-6beb-4c9d-aa2b-0f7804a3e214"

func TestTerradragonCusdalf(t *testing.T) {
	t.Run("keeps only its controller's mana tapped", func(t *testing.T) {
		scn := scenario.New()
		owner := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))

		owner.Player.SpawnCard(terradragonCusdalfUID, match.HAND)
		cusdalf, err := scn.FindCard(owner.Player, match.HAND, terradragonCusdalfUID)
		require.NoError(t, err)
		cusdalf, err = owner.Player.MoveCard(cusdalf.ID, match.HAND, match.BATTLEZONE, "terradragon_cusdalf_test_setup")
		require.NoError(t, err)

		for range 2 {
			owner.Player.SpawnCard(scowlingTomatoUID, match.MANAZONE)
			opponent.Player.SpawnCard(scowlingTomatoUID, match.MANAZONE)
		}
		ownerMana := tapAllTerradragonCusdalfTestMana(t, owner.Player)
		opponentMana := tapAllTerradragonCusdalfTestMana(t, opponent.Player)
		ownerMana[0].Tapped = false
		cusdalf.Tapped = true

		assert.Equal(t, "Terradragon Cusdalf", cusdalf.Name)
		assert.Equal(t, 7000, cusdalf.Power)
		assert.Equal(t, 5, cusdalf.ManaCost)
		assert.Equal(t, []string{civ.Nature}, cusdalf.Civs)
		assert.True(t, cusdalf.HasFamily(family.EarthDragon))

		require.NoError(t, scn.ActionEndTurn(owner))
		for _, manaCard := range opponentMana {
			assert.False(t, manaCard.Tapped, "Cusdalf must not affect the opponent's mana")
		}

		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.False(t, ownerMana[0].Tapped, "mana that was already untapped must stay untapped")
		assert.True(t, ownerMana[1].Tapped, "tapped mana must not untap")
		assert.False(t, cusdalf.Tapped, "Cusdalf itself still untaps in the battle zone")
		assert.True(t, cusdalf.HasCondition(cnd.DoubleBreaker))
		assert.Equal(t, 7000, scn.Match.GetPower(cusdalf, false))
		assert.Equal(t, 11000, scn.Match.GetPower(cusdalf, true))
	})

	t.Run("does not prevent mana untapping from another zone", func(t *testing.T) {
		scn := scenario.New()
		owner := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))

		owner.Player.SpawnCard(terradragonCusdalfUID, match.HAND)
		owner.Player.SpawnCard(scowlingTomatoUID, match.MANAZONE)
		ownerMana := tapAllTerradragonCusdalfTestMana(t, owner.Player)

		require.NoError(t, scn.ActionEndTurn(owner))
		require.NoError(t, scn.ActionEndTurn(opponent))
		for _, manaCard := range ownerMana {
			assert.False(t, manaCard.Tapped)
		}
	})

	t.Run("kept mana is visible to opposing continuous power effects", func(t *testing.T) {
		scn := scenario.New()
		owner := scn.Match.CurrentPlayer()
		opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))

		owner.Player.SpawnCard(terradragonCusdalfUID, match.HAND)
		cusdalf, err := scn.FindCard(owner.Player, match.HAND, terradragonCusdalfUID)
		require.NoError(t, err)
		_, err = owner.Player.MoveCard(cusdalf.ID, match.HAND, match.BATTLEZONE, "terradragon_cusdalf_test_setup")
		require.NoError(t, err)

		opponent.Player.SpawnCard(kingOquanosUID, match.HAND)
		king, err := scn.FindCard(opponent.Player, match.HAND, kingOquanosUID)
		require.NoError(t, err)
		_, err = opponent.Player.MoveCard(king.ID, match.HAND, match.BATTLEZONE, "terradragon_cusdalf_test_setup")
		require.NoError(t, err)

		for range 2 {
			owner.Player.SpawnCard(scowlingTomatoUID, match.MANAZONE)
		}
		tapAllTerradragonCusdalfTestMana(t, owner.Player)

		require.NoError(t, scn.ActionEndTurn(owner))
		require.NoError(t, scn.ActionEndTurn(opponent))

		assert.Equal(t, 6000, scn.Match.GetPower(king, false))
		assert.True(t, king.HasCondition(cnd.DoubleBreaker))
	})
}

func tapAllTerradragonCusdalfTestMana(t *testing.T, player *match.Player) []*match.Card {
	t.Helper()

	mana, err := player.Container(match.MANAZONE)
	require.NoError(t, err)
	for _, manaCard := range mana {
		manaCard.Tapped = true
	}
	return append([]*match.Card(nil), mana...)
}
