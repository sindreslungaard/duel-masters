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
	jagilaUID   = "e69f512c-86a5-4afa-b76f-77e70ed103b6"
	jagilaSetup = "jagila_the_hidden_pillager_test_setup"
)

func TestJagilaTheHiddenPillager(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		jagila := putCardInBattlezone(t, scn, player.Player, jagilaUID, jagilaSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, jagila, "Jagila, the Hidden Pillager", 3000, 5, []string{civ.Darkness})
		assert.True(t, jagila.HasFamily(family.PandorasBox))
		assert.True(t, jagila.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count it strips three cards at random", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)

		jagila := spawnForLater(t, player, jagilaUID)
		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(handBefore), 3)

		putIntoPlay(t, scn, player, jagila)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore)-3)
	})

	t.Run("without the count the hand is untouched", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		jagila := spawnForLater(t, player, jagilaUID)
		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		putIntoPlay(t, scn, player, jagila)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore))
	})
}
