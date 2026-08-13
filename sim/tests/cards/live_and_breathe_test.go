package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	liveAndBreatheUID = "b7d5c565-dbcc-4883-a175-91a6f5b72f25"
	liveSummonUID     = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (fire, cost 2)
	liveOtherUID      = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars
	liveAndBreatheSrc = "live_and_breathe_test_setup"
)

func TestLiveAndBreathe(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		_, player, _ := setupDuel(t)

		spell, err := player.Player.SpawnCard(liveAndBreatheUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, "Live and Breathe", spell.Name)
		assert.Equal(t, 3, spell.ManaCost)
		assert.Equal(t, []string{civ.Light, civ.Nature}, spell.Civs)
	})

	t.Run("summoning a creature fetches a copy of it into the battle zone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// A copy waiting in the deck, below the top so the draw step cannot
		// take it first.
		player.Player.DestroyDeck()
		for range 3 {
			_, err := player.Player.SpawnCard(liveOtherUID, match.DECK)
			require.NoError(t, err)
		}
		copyInDeck, err := player.Player.SpawnCard(liveSummonUID, match.DECK)
		require.NoError(t, err)
		for range 4 {
			_, err := player.Player.SpawnCard(liveOtherUID, match.DECK)
			require.NoError(t, err)
		}

		toSummon := spawnForLater(t, player, liveSummonUID)

		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, liveAndBreatheUID)
		require.NoError(t, scn.WaitForEventLoop())

		// Seeded after the cast: Live and Breathe is multicolored, so paying
		// for it can legitimately consume mana of any civilization, and the
		// creature needs its own fire left over.
		for range 2 {
			_, err := player.Player.SpawnCard(liveSummonUID, match.MANAZONE)
			require.NoError(t, err)
		}

		require.NoError(t, scn.ActionPlayCard(player, toSummon.ID))
		answerInTurn(t, scn, player, copyInDeck.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, toSummon.Zone)
		assert.Equal(t, match.BATTLEZONE, copyInDeck.Zone, "the copy joins it")
	})

	t.Run("the fetched creature does not trigger another search", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		player.Player.DestroyDeck()
		for range 3 {
			_, err := player.Player.SpawnCard(liveOtherUID, match.DECK)
			require.NoError(t, err)
		}
		first, err := player.Player.SpawnCard(liveSummonUID, match.DECK)
		require.NoError(t, err)
		second, err := player.Player.SpawnCard(liveSummonUID, match.DECK)
		require.NoError(t, err)
		for range 4 {
			_, err := player.Player.SpawnCard(liveOtherUID, match.DECK)
			require.NoError(t, err)
		}

		toSummon := spawnForLater(t, player, liveSummonUID)

		passTurnToSelf(t, scn, player, opponent)

		castSpell(t, scn, player, liveAndBreatheUID)
		require.NoError(t, scn.WaitForEventLoop())

		// Seeded after the cast: Live and Breathe is multicolored, so paying
		// for it can legitimately consume mana of any civilization, and the
		// creature needs its own fire left over.
		for range 2 {
			_, err := player.Player.SpawnCard(liveSummonUID, match.MANAZONE)
			require.NoError(t, err)
		}

		require.NoError(t, scn.ActionPlayCard(player, toSummon.ID))
		answerInTurn(t, scn, player, first.ID)
		require.NoError(t, scn.WaitForEventLoop())

		// A creature put into play by an effect is not summoned, so the second
		// copy stays where it is rather than the search running forever.
		assert.Equal(t, match.BATTLEZONE, first.Zone)
		assert.NotEqual(t, match.BATTLEZONE, second.Zone)
	})

	t.Run("the effect is over on the following turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		player.Player.DestroyDeck()
		for range 3 {
			_, err := player.Player.SpawnCard(liveOtherUID, match.DECK)
			require.NoError(t, err)
		}
		copyInDeck, err := player.Player.SpawnCard(liveSummonUID, match.DECK)
		require.NoError(t, err)
		for range 6 {
			_, err := player.Player.SpawnCard(liveOtherUID, match.DECK)
			require.NoError(t, err)
		}

		toSummon := spawnForLater(t, player, liveSummonUID)

		passTurnToSelf(t, scn, player, opponent)
		castSpell(t, scn, player, liveAndBreatheUID)
		require.NoError(t, scn.WaitForEventLoop())

		for range 2 {
			_, err := player.Player.SpawnCard(liveSummonUID, match.MANAZONE)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.ActionPlayCard(player, toSummon.ID))
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment should ask anything")
		assert.NotEqual(t, match.BATTLEZONE, copyInDeck.Zone)
	})
}
