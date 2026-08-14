package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	franticChieftainUID       = "8a362a6e-8fd8-41e7-ad6c-bf27be952698"
	franticChieftainCheapUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (fire, cost 2)
	franticChieftainCostlyUID = "aeaaf98d-938f-46d1-a271-49a86f668ae6" // Typhoon Crawler (cost 6)
	franticChieftainSetupSrc  = "frantic_chieftain_test_setup"
)

func TestFranticChieftain(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// Kept in hand: putting it into play makes it bounce itself, since it
		// costs 2 and would be the only legal target for its own effect.
		card := spawnForLater(t, player, franticChieftainUID)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Frantic Chieftain", 2000, 2, []string{civ.Water})
		assert.True(t, card.HasFamily(family.Merfolk))
	})

	t.Run("it returns a cheap creature to hand", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		cheap := putCardInBattlezone(t, scn, player.Player, franticChieftainCheapUID, franticChieftainSetupSrc)

		inPlay := summonWithOwnMana(t, scn, player, franticChieftainUID)

		// Two legal choices, the cheap creature and the chieftain itself. The
		// mana it was paid with is also a Frantic Chieftain, but that is in the
		// mana zone rather than the battle zone.
		answerInTurn(t, scn, player, cheap.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, cheap.Zone)
		assert.Equal(t, match.BATTLEZONE, inPlay.Zone)
	})

	t.Run("it can be forced to return itself", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		chieftain := summonWithOwnMana(t, scn, player, franticChieftainUID)
		require.NoError(t, scn.WaitForEventLoop())

		// It costs 2, so with nothing else on the board it is the only legal
		// choice and the mandatory return takes it.
		assert.Equal(t, match.HAND, chieftain.Zone)
	})

	t.Run("expensive creatures are not eligible", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		expensive := putCardInBattlezone(t, scn, player.Player, franticChieftainCostlyUID, franticChieftainSetupSrc)

		chieftain := summonWithOwnMana(t, scn, player, franticChieftainUID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, expensive.Zone, "6 is more than 4")
		assert.Equal(t, match.HAND, chieftain.Zone, "so the chieftain returns itself")
	})
}
