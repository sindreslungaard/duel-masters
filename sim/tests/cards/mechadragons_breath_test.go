package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mechadragonsBreathUID = "6dcc11be-ceab-4a72-8d00-937b1f43bbd6"
	mechadragonTwoKUID    = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	mechadragonFourKUID   = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	mechadragonSetupSrc   = "mechadragons_breath_test_setup"
)

func TestMechadragonsBreath(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(mechadragonsBreathUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Mechadragon's Breath", spell.Name)
		assert.Equal(t, 6, spell.ManaCost)
		assert.Equal(t, []string{civ.Fire}, spell.Civs)
		assert.Equal(t, []string{civ.Fire}, spell.ManaRequirement)
	})

	t.Run("only creatures of exactly that power are destroyed", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		mine := putCardInBattlezone(t, scn, player.Player, mechadragonTwoKUID, mechadragonSetupSrc)
		theirSmall := putCardInBattlezone(t, scn, opponent.Player, mechadragonTwoKUID, mechadragonSetupSrc)
		theirBig := putCardInBattlezone(t, scn, opponent.Player, mechadragonFourKUID, mechadragonSetupSrc)

		castSpell(t, scn, player, mechadragonsBreathUID)

		require.NoError(t, scn.SubmitCount(player, 2000))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, mine.Zone, "its caster is not spared")
		assert.Equal(t, match.GRAVEYARD, theirSmall.Zone)
		assert.Equal(t, match.BATTLEZONE, theirBig.Zone, "4000 is not 2000")
	})

	t.Run("a power nobody has costs nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirs := putCardInBattlezone(t, scn, opponent.Player, mechadragonTwoKUID, mechadragonSetupSrc)

		castSpell(t, scn, player, mechadragonsBreathUID)

		require.NoError(t, scn.SubmitCount(player, 5000))
		settleTurn(t, scn)

		assert.Equal(t, match.BATTLEZONE, theirs.Zone)
	})

	t.Run("the number is refused above its ceiling", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirs := putCardInBattlezone(t, scn, opponent.Player, mechadragonFourKUID, mechadragonSetupSrc)

		castSpell(t, scn, player, mechadragonsBreathUID)

		warningStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.SubmitCount(player, 7000))
		require.NoError(t, scn.WaitForMessage(player, warningStart, "action_error"))

		// The prompt is still open, so a legal number can still be given.
		require.NoError(t, scn.SubmitCount(player, 4000))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, theirs.Zone)
	})

	t.Run("effective power counts, not printed power", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		buffed := putCardInBattlezone(t, scn, opponent.Player, mechadragonTwoKUID, mechadragonSetupSrc)
		buffed.AddUniqueSourceCondition(cnd.PowerAmplifier, 1000, mechadragonSetupSrc)

		castSpell(t, scn, player, mechadragonsBreathUID)

		require.NoError(t, scn.SubmitCount(player, 3000))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, buffed.Zone, "2000 printed, 3000 in play")
	})
}
