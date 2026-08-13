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
	elfXUID          = "60d8c6a6-20c1-425c-9ecc-b56981a70e21"
	elfXCreatureUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5"
	elfXNatureManaID = "1d72eb3e-5185-449a-a16f-391bd2338343"
)

func TestElfX(t *testing.T) {
	t.Run("applies its reduction in the turn it is summoned", func(t *testing.T) {
		scn := scenario.New()
		owner := scn.Match.CurrentPlayer()

		elfX, err := owner.Player.SpawnCard(elfXUID, match.HAND)
		require.NoError(t, err)
		require.NotNil(t, elfX)

		for range 4 {
			_, err := owner.Player.SpawnCard(elfXNatureManaID, match.MANAZONE)
			require.NoError(t, err)
		}
		_, err = owner.Player.SpawnCard(elfXCreatureUID, match.MANAZONE)
		require.NoError(t, err)

		creature, err := scn.FindCard(owner.Player, match.HAND, elfXCreatureUID)
		require.NoError(t, err)
		require.True(t, creature.HasCondition(cnd.Creature))
		assert.False(t, creature.HasCondition(cnd.ReducedCost))

		assert.Equal(t, "Elf-X", elfX.Name)
		assert.Equal(t, 2000, elfX.Power)
		assert.Equal(t, 4, elfX.ManaCost)
		assert.Equal(t, []string{civ.Nature}, elfX.Civs)
		assert.Equal(t, []string{civ.Nature}, elfX.ManaRequirement)
		assert.True(t, elfX.HasFamily(family.TreeFolk))

		require.NoError(t, scn.ActionPlayCard(owner, elfX.ID))
		assert.True(t, creature.HasCondition(cnd.ReducedCost))

		mana, err := owner.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.True(t, owner.Player.CanPlayCard(creature, mana), "the remaining fire mana should pay the reduced cost")

		require.NoError(t, scn.ActionPlayCard(owner, creature.ID))
		assert.Equal(t, match.BATTLEZONE, creature.Zone)
	})

	t.Run("multiple copies stack and clean up only their own reduction", func(t *testing.T) {
		scn := scenario.New()
		owner := scn.Match.CurrentPlayer()
		creature, err := scn.FindCard(owner.Player, match.HAND, elfXCreatureUID)
		require.NoError(t, err)

		first := moveElfXToBattleZone(t, owner.Player)
		second := moveElfXToBattleZone(t, owner.Player)
		assert.Equal(t, 2, reducedCostSourceCount(creature))

		moved, err := owner.Player.MoveCard(first.ID, match.BATTLEZONE, match.GRAVEYARD, "elf_x_test")
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
		assert.Equal(t, 1, reducedCostSourceCount(creature))

		moved, err = owner.Player.MoveCard(second.ID, match.BATTLEZONE, match.GRAVEYARD, "elf_x_test")
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
		assert.Zero(t, reducedCostSourceCount(creature))
	})
}

func moveElfXToBattleZone(t *testing.T, player *match.Player) *match.Card {
	t.Helper()

	elfX, err := player.SpawnCard(elfXUID, match.HAND)
	require.NoError(t, err)
	elfX, err = player.MoveCard(elfX.ID, match.HAND, match.BATTLEZONE, "elf_x_test")
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, elfX.Zone)

	return elfX
}

func reducedCostSourceCount(card *match.Card) int {
	count := 0
	for _, condition := range card.Conditions() {
		if condition.ID == cnd.ReducedCost {
			count++
		}
	}

	return count
}
