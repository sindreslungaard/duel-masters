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
	riseAndShineUID     = "47abe92b-a677-4c6a-b9da-033b28bcf374"
	riseAndShineBlocker = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (blocker)
	riseAndShinePlain   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	riseAndShineSrc     = "rise_and_shine_test_setup"
)

func TestRiseAndShine(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		spell, err := player.Player.SpawnCard(riseAndShineUID, match.HAND)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Rise and Shine", spell.Name)
		assert.Equal(t, 4, spell.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Water}, spell.Civs)
		assert.True(t, spell.IsMulticolored())
		assert.True(t, spell.HasCondition(cnd.ShieldTrigger))
	})

	t.Run("it enters the mana zone tapped", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		card, err := player.Player.SpawnCard(riseAndShineUID, match.HAND)
		require.NoError(t, err)

		moved, err := player.Player.MoveCard(card.ID, match.HAND, match.MANAZONE, riseAndShineSrc)
		require.NoError(t, err)

		assert.True(t, moved.Tapped)
	})

	t.Run("it takes a blocker and buries the rest", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// Seeded below the top so the draw step cannot take the blocker first.
		player.Player.DestroyDeck()
		for range 2 {
			_, err := player.Player.SpawnCard(riseAndShinePlain, match.DECK)
			require.NoError(t, err)
		}
		blocker, err := player.Player.SpawnCard(riseAndShineBlocker, match.DECK)
		require.NoError(t, err)
		for range 6 {
			_, err := player.Player.SpawnCard(riseAndShinePlain, match.DECK)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		deckBefore, err := player.Player.Container(match.DECK)
		require.NoError(t, err)

		castSpell(t, scn, player, riseAndShineUID)

		// One blocker among the four, so it is taken without asking. The three
		// left over then have to be put in order for the bottom of the deck.
		answerOrderPrompt(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, blocker.Zone)

		deck, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		assert.Len(t, deck, len(deckBefore)-1)
	})

	t.Run("no blocker in the four means nothing is taken", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		player.Player.DestroyDeck()
		for range 9 {
			_, err := player.Player.SpawnCard(riseAndShinePlain, match.DECK)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		deckBefore, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		castSpell(t, scn, player, riseAndShineUID)

		// Nothing was taken, so all four are ordered onto the bottom.
		answerOrderPrompt(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		deck, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		// The spell left the hand and nothing came back; the deck is untouched
		// beyond being reordered.
		assert.Len(t, deck, len(deckBefore))
		assert.Len(t, hand, len(handBefore))
	})
}
