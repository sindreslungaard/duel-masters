package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	gandavalsStaplerUID   = "d0429aea-e222-482e-a23a-49db1df89c98"
	gandavalsArrivalUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	gandavalsStaplerSetup = "gandavals_stapler_test_setup"
)

func TestGandavalsStapler(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, gandavalsStaplerUID, gandavalsStaplerSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Gandaval's Stapler", 3000, 2, []string{civ.Fire})
		assert.True(t, card.HasFamily(family.Xenoparts))
	})

	t.Run("another creature arriving taps it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		stapler := putCardInBattlezone(t, scn, player.Player, gandavalsStaplerUID, gandavalsStaplerSetup)
		passTurnToSelf(t, scn, player, opponent)
		require.False(t, stapler.Tapped)

		putCardInBattlezone(t, scn, opponent.Player, gandavalsArrivalUID, gandavalsStaplerSetup)
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, stapler.Tapped)
	})
}
