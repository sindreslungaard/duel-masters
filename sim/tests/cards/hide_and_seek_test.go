package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	hideAndSeekUID      = "2072f6de-c78c-492c-b026-4d94b148e8a2"
	hideAndSeekSetupSrc = "hide_and_seek_test_setup"
)

func TestHideAndSeek(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(hideAndSeekUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Hide and Seek", spell.Name)
		assert.Equal(t, 4, spell.ManaCost)
		assert.Equal(t, []string{civ.Water, civ.Darkness}, spell.Civs)
	})

	t.Run("it bounces a creature and then takes a card at random", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirs := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, hideAndSeekSetupSrc)
		ours := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, hideAndSeekSetupSrc)

		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		castSpell(t, scn, player, hideAndSeekUID)
		require.NoError(t, scn.WaitForEventLoop())

		// One legal target, so it is taken without asking. The bounced creature
		// joins the hand before the discard, which makes it a candidate for the
		// discard itself, so only its departure from the battle zone is certain.
		assert.NotEqual(t, match.BATTLEZONE, theirs.Zone)
		assert.Equal(t, match.BATTLEZONE, ours.Zone, "only the opponent's creatures are eligible")

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore), "one card in, one card out")
	})

	t.Run("the discard happens even with no creature to bounce", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		handBefore, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)

		castSpell(t, scn, player, hideAndSeekUID)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore)-1, "the discard is a sentence of its own")
	})
}
