package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	copperLocustUID = "2b5cb952-39ea-4264-8f73-92b912102021"
	// Armored Cannon Balbaro evolves from a Human, and the scenario deck is made
	// of Immortal Baron, Vorg, which is one.
	copperLocustEvolutionUID = "24353d06-89ef-4867-9513-485750d01e10"
	copperLocustBaseUID      = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5"
	// Ultra Mantis, Scourge of Fate evolves from a Giant Insect, which is what
	// Copper Locust is.
	copperLocustInsectEvolutionUID = "4269337e-5772-4d22-9474-e3068cf21de0"
	copperLocustSrc                = "copper_locust_test_setup"
)

func TestCopperLocust(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, copperLocustUID, copperLocustSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Copper Locust", 5000, 3, []string{civ.Nature})
		assert.True(t, card.HasFamily(family.GiantInsect))
	})

	t.Run("its controller evolving kills it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		locust := putCardInBattlezone(t, scn, player.Player, copperLocustUID, copperLocustSrc)
		base := putCardInBattlezone(t, scn, player.Player, copperLocustBaseUID, copperLocustSrc)

		passTurnToSelf(t, scn, player, opponent)

		// One legal base, so the evolution takes it without asking.
		evolution := summonWithOwnMana(t, scn, player, copperLocustEvolutionUID)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, locust.Zone)
		assert.Equal(t, match.BATTLEZONE, evolution.Zone)
		assert.Equal(t, match.HIDDENZONE, base.Zone)
	})

	t.Run("the opponent evolving kills it just the same", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		locust := putCardInBattlezone(t, scn, player.Player, copperLocustUID, copperLocustSrc)
		putCardInBattlezone(t, scn, opponent.Player, copperLocustBaseUID, copperLocustSrc)

		require.NoError(t, scn.ActionEndTurn(player))

		evolution := summonWithOwnMana(t, scn, opponent, copperLocustEvolutionUID)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, locust.Zone, "\"a player\" is either player")
		assert.Equal(t, match.BATTLEZONE, evolution.Zone)
	})

	t.Run("an ordinary summon leaves it alone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		locust := putCardInBattlezone(t, scn, player.Player, copperLocustUID, copperLocustSrc)

		passTurnToSelf(t, scn, player, opponent)

		summoned := summonWithOwnMana(t, scn, player, copperLocustBaseUID)
		settleTurn(t, scn)

		assert.Equal(t, match.BATTLEZONE, locust.Zone)
		assert.Equal(t, match.BATTLEZONE, summoned.Zone)
	})

	t.Run("being the base itself costs it nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		locust := putCardInBattlezone(t, scn, player.Player, copperLocustUID, copperLocustSrc)

		passTurnToSelf(t, scn, player, opponent)

		// Ultra Mantis evolves from a Giant Insect, and the locust is the only
		// one in play, so it is what gets evolved on.
		evolution := summonWithOwnMana(t, scn, player, copperLocustInsectEvolutionUID)
		settleTurn(t, scn)

		// It is under the evolution rather than in the battle zone by the time
		// it would destroy itself, so it goes nowhere.
		assert.Equal(t, match.HIDDENZONE, locust.Zone)
		assert.Equal(t, match.BATTLEZONE, evolution.Zone)
	})
}
