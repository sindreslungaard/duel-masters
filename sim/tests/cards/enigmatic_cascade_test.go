package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const enigmaticCascadeUID = "2e9096a9-dd85-46f4-9159-aa73abe30165"

func TestEnigmaticCascade(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(enigmaticCascadeUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Enigmatic Cascade", spell.Name)
		assert.Equal(t, 4, spell.ManaCost)
		assert.Equal(t, []string{civ.Water}, spell.Civs)
		assert.Equal(t, []string{civ.Water}, spell.ManaRequirement)
	})

	t.Run("what is discarded is drawn back", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(handBefore), 2)
		handCount := len(handBefore)

		first, second := handBefore[0], handBefore[1]

		castSpell(t, scn, player, enigmaticCascadeUID)

		answerInTurn(t, scn, player, first.ID, second.ID)
		settleTurn(t, scn)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		assert.Equal(t, match.GRAVEYARD, first.Zone)
		assert.Equal(t, match.GRAVEYARD, second.Zone)
		// The spell was added to this hand and then cast, so two out and two in
		// leaves it exactly as it started.
		assert.Len(t, hand, handCount)
	})

	t.Run("any number includes none", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(handBefore)

		graveyardBefore, err := player.Player.Container(match.GRAVEYARD)
		require.NoError(t, err)
		graveyardCount := len(graveyardBefore)

		spell := castSpell(t, scn, player, enigmaticCascadeUID)

		cancelInTurn(t, scn, player)
		settleTurn(t, scn)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		graveyard, err := player.Player.Container(match.GRAVEYARD)
		require.NoError(t, err)

		assert.Len(t, hand, handCount, "nothing was discarded and nothing was drawn")
		assert.Len(t, graveyard, graveyardCount+1, "only the spell arrived")
		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})

	t.Run("an empty hand asks nothing", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(enigmaticCascadeUID, match.HAND)
		require.NoError(t, err)

		for range spell.ManaCost {
			_, err := player.Player.SpawnCard(enigmaticCascadeUID, match.MANAZONE)
			require.NoError(t, err)
		}

		// Emptied down to the spell, which is gone from the hand by the time its
		// own effect resolves. Snapshotted first, because moving cards out of a
		// zone rewrites the slice being ranged over.
		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		snapshot := make([]*match.Card, len(hand))
		copy(snapshot, hand)

		for _, card := range snapshot {
			if card.ID == spell.ID {
				continue
			}

			_, err := player.Player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, sharedTestSrc)
			require.NoError(t, err)
		}

		require.NoError(t, scn.ActionPlayCard(player, spell.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, spell.Zone)
	})
}
