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
	klujadrasUID   = "89eddde4-a2fe-45c3-ad47-813c87ac6d37"
	klujadrasSetup = "klujadras_test_setup"
)

func TestKlujadras(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		klujadras := putCardInBattlezone(t, scn, player.Player, klujadrasUID, klujadrasSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, klujadras, "Klujadras", 4000, 7, []string{civ.Water})
		assert.True(t, klujadras.HasFamily(family.SeaHacker))
		assert.True(t, klujadras.HasCondition(cnd.WaveStriker))
	})

	t.Run("each player draws for their own wave strikers", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)
		putCardInBattlezone(t, scn, opponent.Player, waveStrikerFillerUID, klujadrasSetup)
		putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, klujadrasSetup)

		klujadras := spawnForLater(t, player, klujadrasUID)
		passTurnToSelf(t, scn, player, opponent)

		myHandBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		theirHandBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		putIntoPlay(t, scn, player, klujadras)

		myHand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		theirHand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		// Three on the caster's side once Klujadras itself has arrived, one on
		// the opponent's; the plain creature does not count.
		assert.Len(t, myHand, len(myHandBefore)-1+3, "Klujadras left the hand and drew three")
		assert.Len(t, theirHand, len(theirHandBefore)+1)
	})
}
