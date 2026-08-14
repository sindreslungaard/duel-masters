package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	whirlingMalianUID   = "a328f1ed-12b9-467e-b621-19f676e75714"
	whirlingArrivalUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	whirlingMalianSetup = "whirling_warrior_malian_test_setup"
)

func TestWhirlingWarriorMalian(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, whirlingMalianUID, whirlingMalianSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Whirling Warrior Malian", 6000, 4, []string{civ.Fire})
		assert.True(t, card.HasFamily(family.Armorloid))
		assert.False(t, card.Tapped, "nothing has arrived yet")
	})

	t.Run("another creature arriving taps it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		malian := putCardInBattlezone(t, scn, player.Player, whirlingMalianUID, whirlingMalianSetup)
		passTurnToSelf(t, scn, player, opponent)
		require.False(t, malian.Tapped)

		putCardInBattlezone(t, scn, player.Player, whirlingArrivalUID, whirlingMalianSetup)
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, malian.Tapped)
	})

	t.Run("the opponent's creatures count too", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		malian := putCardInBattlezone(t, scn, player.Player, whirlingMalianUID, whirlingMalianSetup)
		passTurnToSelf(t, scn, player, opponent)

		putCardInBattlezone(t, scn, opponent.Player, whirlingArrivalUID, whirlingMalianSetup)
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, malian.Tapped, "the printed text is \"another creature\", not \"another of yours\"")
	})
}
