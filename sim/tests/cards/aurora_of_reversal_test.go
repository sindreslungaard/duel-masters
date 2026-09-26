package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	auroraOfReversalUID       = "9cff685f-3491-400b-b9d4-43c383193dc7"
	auroraOfReversalFillerUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
)

func TestAuroraOfReversal(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(auroraOfReversalUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Aurora of Reversal", spell.Name)
		assert.Equal(t, 5, spell.ManaCost)
		assert.Equal(t, []string{civ.Nature}, spell.Civs)
	})

	t.Run("the cap on how many shields to move is based on its own controller's shields, not the opponent's", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// Give the caster more shields than the opponent, so a cap mistakenly
		// based on the opponent's (smaller) count would be observable.
		for range 2 {
			_, err := player.Player.SpawnCard(auroraOfReversalFillerUID, match.SHIELDZONE)
			require.NoError(t, err)
		}

		myShields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		myShieldCount := len(myShields)

		oppShields, err := opponent.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		require.Less(t, len(oppShields), myShieldCount, "the opponent must have fewer shields for this to test anything")

		spell := castSpell(t, scn, player, auroraOfReversalUID)

		action, err := scn.LatestAction(player, 0)
		require.NoError(t, err, "expected the shield-selection prompt")
		assert.True(t, action.Cancellable, "declining is how zero shields is chosen")
		assert.Equal(t, myShieldCount, action.MaxSelections, "capped by its own controller's shields, not the opponent's")

		ids := make([]string, 0, len(action.Cards))
		for _, c := range action.Cards {
			ids = append(ids, c.CardID)
		}
		require.Len(t, ids, myShieldCount, "every one of its controller's own shields is offered")
		answerInTurn(t, scn, player, ids...)

		remainingShields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Empty(t, remainingShields, "every shield selected was moved")

		mana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, mana, spell.ManaCost+myShieldCount, "the mana paid for the spell plus every shield moved")
	})

	t.Run("it may be declined, moving nothing", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		shieldsBefore, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		shieldCount := len(shieldsBefore)

		castSpell(t, scn, player, auroraOfReversalUID)

		_, err = scn.LatestAction(player, 0)
		require.NoError(t, err, "expected the shield-selection prompt")

		cancelInTurn(t, scn, player)

		remainingShields, err := player.Player.Container(match.SHIELDZONE)
		require.NoError(t, err)
		assert.Len(t, remainingShields, shieldCount, "nothing was moved")
	})
}
