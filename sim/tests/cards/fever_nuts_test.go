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
	feverNutsUID      = "5c9e07c1-681c-41f6-975e-0608ff90ab6b"
	feverNutsOtherUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000, cost 2)
	feverNutsSetupSrc = "fever_nuts_test_setup"
)

func TestFeverNuts(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, feverNutsUID, feverNutsSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Fever Nuts", 1000, 3, []string{civ.Nature})
		assert.True(t, card.HasFamily(family.WildVeggies))
	})

	t.Run("creatures on both sides get the discount", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		mine, err := player.Player.SpawnCard(feverNutsOtherUID, match.HAND)
		require.NoError(t, err)
		theirs, err := opponent.Player.SpawnCard(feverNutsOtherUID, match.HAND)
		require.NoError(t, err)

		// Installed before the turn passes, because a persistent effect only
		// starts running on the next event after it is applied.
		putCardInBattlezone(t, scn, player.Player, feverNutsUID, feverNutsSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.True(t, mine.HasCondition(cnd.ReducedCost))
		assert.True(t, theirs.HasCondition(cnd.ReducedCost), "the printed text helps the opponent too")
	})

	t.Run("the discount is really paid", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		creature, err := player.Player.SpawnCard(feverNutsOtherUID, match.HAND)
		require.NoError(t, err)

		// One mana short of the printed cost of 2.
		_, err = player.Player.SpawnCard(feverNutsOtherUID, match.MANAZONE)
		require.NoError(t, err)

		putCardInBattlezone(t, scn, player.Player, feverNutsUID, feverNutsSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionPlayCard(player, creature.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.BATTLEZONE, creature.Zone, "a 2 cost creature summoned off one mana")
	})

	t.Run("the discount leaves with it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		mine, err := player.Player.SpawnCard(feverNutsOtherUID, match.HAND)
		require.NoError(t, err)

		nuts := putCardInBattlezone(t, scn, player.Player, feverNutsUID, feverNutsSetupSrc)
		passTurnToSelf(t, scn, player, opponent)
		require.True(t, mine.HasCondition(cnd.ReducedCost))

		_, err = player.Player.MoveCard(nuts.ID, match.BATTLEZONE, match.GRAVEYARD, feverNutsSetupSrc)
		require.NoError(t, err)
		settleTurn(t, scn)

		assert.False(t, mine.HasCondition(cnd.ReducedCost))
	})

	t.Run("two copies do not stack past the floor", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		creature, err := player.Player.SpawnCard(feverNutsOtherUID, match.HAND)
		require.NoError(t, err)

		putCardInBattlezone(t, scn, player.Player, feverNutsUID, feverNutsSetupSrc)
		putCardInBattlezone(t, scn, player.Player, feverNutsUID, feverNutsSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		conditions := 0
		for _, condition := range creature.Conditions() {
			if condition.ID == cnd.ReducedCost {
				conditions++
			}
		}

		// Each copy contributes its own reduction, and fx.Creature keeps the
		// total from dropping below one mana when it works out what to pay.
		assert.Equal(t, 2, conditions)

		_, err = player.Player.SpawnCard(feverNutsOtherUID, match.MANAZONE)
		require.NoError(t, err)

		require.NoError(t, scn.ActionPlayCard(player, creature.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.BATTLEZONE, creature.Zone)
	})
}
