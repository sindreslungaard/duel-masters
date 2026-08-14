package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	clonedNightmareUID = "060ee9a8-f238-41f3-9e46-3b467c487a5e"
	clonedNightmareSrc = "cloned_nightmare_test_setup"
)

func TestClonedNightmare(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(clonedNightmareUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Cloned Nightmare", spell.Name)
		assert.Equal(t, 3, spell.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, spell.Civs)
		assert.Equal(t, []string{civ.Darkness}, spell.ManaRequirement)
	})

	t.Run("with no copies buried it takes one card and asks nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(hand)

		castSpell(t, scn, player, clonedNightmareUID)
		settleTurn(t, scn)

		hand, err = opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount-1)
	})

	t.Run("a buried copy offers one more", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		_, err := player.Player.SpawnCard(clonedNightmareUID, match.GRAVEYARD)
		require.NoError(t, err)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(hand)

		castSpell(t, scn, player, clonedNightmareUID)

		require.NoError(t, scn.SubmitCount(player, 1))
		settleTurn(t, scn)

		hand, err = opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount-2)
	})

	t.Run("the extra card may be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		_, err := opponent.Player.SpawnCard(clonedNightmareUID, match.GRAVEYARD)
		require.NoError(t, err)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(hand)

		castSpell(t, scn, player, clonedNightmareUID)

		require.NoError(t, scn.SubmitCount(player, 0))
		settleTurn(t, scn)

		hand, err = opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount-1, "the first card is never optional")
	})

	t.Run("more than the copies allow is refused", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		_, err := player.Player.SpawnCard(clonedNightmareUID, match.GRAVEYARD)
		require.NoError(t, err)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(hand)

		castSpell(t, scn, player, clonedNightmareUID)

		warningStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.SubmitCount(player, 2))
		require.NoError(t, scn.WaitForMessage(player, warningStart, "action_error"))

		require.NoError(t, scn.SubmitCount(player, 1))
		settleTurn(t, scn)

		hand, err = opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount-2)
	})

	t.Run("an empty hand asks nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		_, err := player.Player.SpawnCard(clonedNightmareUID, match.GRAVEYARD)
		require.NoError(t, err)

		emptyHand(t, opponent, clonedNightmareSrc)

		castSpell(t, scn, player, clonedNightmareUID)
		settleTurn(t, scn)

		hand, err := opponent.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Empty(t, hand)
	})
}
