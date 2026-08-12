package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	miraculousSnareUID      = "3d4863c7-9586-4e9a-af76-9ce1b5b332e2"
	miraculousSnareSetupSrc = "miraculous_snare_test_setup"
)

func TestMiraculousSnare(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(miraculousSnareUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Miraculous Snare", spell.Name)
		assert.Equal(t, 3, spell.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Water}, spell.Civs)
		assert.True(t, spell.IsMulticolored())
	})

	t.Run("the chosen creature becomes one of its owner's shields", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirs := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, miraculousSnareSetupSrc)
		ours := putCardInBattlezone(t, scn, player.Player, immortalBaronVorgUID, miraculousSnareSetupSrc)

		shieldsBefore, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)

		castSpell(t, scn, player, miraculousSnareUID)
		answerInTurn(t, scn, player, theirs.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.SHIELDZONE, theirs.Zone)
		assert.Equal(t, opponent.Player, theirs.Player, "it becomes a shield for its owner, not the caster")
		assert.Equal(t, match.BATTLEZONE, ours.Zone)

		shields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, shields, len(shieldsBefore)+1)
	})

	t.Run("an empty board opens no prompt", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		spell := castSpell(t, scn, player, miraculousSnareUID)
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment should ask anything")
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})
}
