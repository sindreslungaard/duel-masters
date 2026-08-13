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
	clonedDeflectorUID    = "192b1952-e079-4fe2-be03-25d6b655c044"
	clonedDeflectorFoeUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	clonedDeflectorSrc    = "cloned_deflector_test_setup"
)

func TestClonedDeflector(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		spell, err := player.Player.SpawnCard(clonedDeflectorUID, match.HAND)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Cloned Deflector", spell.Name)
		assert.Equal(t, 3, spell.ManaCost)
		assert.Equal(t, []string{civ.Light}, spell.Civs)
		assert.Equal(t, []string{civ.Light}, spell.ManaRequirement)
		assert.True(t, spell.HasCondition(cnd.ShieldTrigger))
	})

	t.Run("with no copies buried it taps exactly one", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		first := putCardInBattlezone(t, scn, opponent.Player, clonedDeflectorFoeUID, clonedDeflectorSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, clonedDeflectorFoeUID, clonedDeflectorSrc)

		castSpell(t, scn, player, clonedDeflectorUID)

		answerInTurn(t, scn, player, first.ID)
		settleTurn(t, scn)

		assert.True(t, first.Tapped)
		assert.False(t, second.Tapped)
	})

	t.Run("each buried copy buys another target", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		first := putCardInBattlezone(t, scn, opponent.Player, clonedDeflectorFoeUID, clonedDeflectorSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, clonedDeflectorFoeUID, clonedDeflectorSrc)

		_, err := opponent.Player.SpawnCard(clonedDeflectorUID, match.GRAVEYARD)
		require.NoError(t, err)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		castSpell(t, scn, player, clonedDeflectorUID)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, action.MinSelections)
		assert.Equal(t, 2, action.MaxSelections, "a copy in the opponent's graveyard counts too")

		answerInTurn(t, scn, player, first.ID, second.ID)
		settleTurn(t, scn)

		assert.True(t, first.Tapped)
		assert.True(t, second.Tapped)
	})

	t.Run("it says in the chat what it tapped", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		target := putCardInBattlezone(t, scn, opponent.Player, clonedDeflectorFoeUID, clonedDeflectorSrc)

		chatStart, err := scn.MessageCount(opponent)
		require.NoError(t, err)

		castSpell(t, scn, player, clonedDeflectorUID)
		settleTurn(t, scn)

		require.True(t, target.Tapped)

		messages, err := scn.ChatMessages(opponent, chatStart)
		require.NoError(t, err)
		assert.Contains(t, messages, "Immortal Baron, Vorg was tapped by Cloned Deflector")
	})

	t.Run("its caster's own creatures are never targets", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		mine := putCardInBattlezone(t, scn, player.Player, clonedDeflectorFoeUID, clonedDeflectorSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, clonedDeflectorFoeUID, clonedDeflectorSrc)

		castSpell(t, scn, player, clonedDeflectorUID)
		settleTurn(t, scn)

		assert.True(t, theirs.Tapped)
		assert.False(t, mine.Tapped)
	})

	t.Run("an empty board asks nothing", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		castSpell(t, scn, player, clonedDeflectorUID)
		settleTurn(t, scn)
	})
}
