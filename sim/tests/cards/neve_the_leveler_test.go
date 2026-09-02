package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	neveTheLevelerUID               = "04965395-67af-4e3e-9b46-1a19efc4e7e8"
	neveTheLevelerFillerCreatureUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	neveTheLevelerSpellUID          = "5883180e-d88c-4f24-b17c-f5a837420147" // Terror Pit
	neveTheLevelerSetupSrc          = "neve_the_leveler_test_setup"
)

func TestNeveTheLeveler(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		neve := putCardInBattlezone(t, scn, player.Player, neveTheLevelerUID, neveTheLevelerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, neve, "Neve, the Leveler", 4000, 6, []string{civ.Nature})
		assert.True(t, neve.HasFamily(family.SnowFaerie))
	})

	t.Run("equal creature counts leave the search unopened", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, opponent.Player, neveTheLevelerFillerCreatureUID, neveTheLevelerSetupSrc)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		// Neve itself is the only creature I have once it resolves, matching the
		// single creature the opponent already has: not "more", so nothing fires.
		summonWithOwnMana(t, scn, player, neveTheLevelerUID)
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment should ask anything")
	})

	t.Run("fewer opponent creatures leave the search unopened", func(t *testing.T) {
		scn, player, _ := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, neveTheLevelerFillerCreatureUID, neveTheLevelerSetupSrc)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, neveTheLevelerUID)
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment should ask anything")
	})

	t.Run("one extra opponent creature offers a search for a single creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, opponent.Player, neveTheLevelerFillerCreatureUID, neveTheLevelerSetupSrc)
		putCardInBattlezone(t, scn, opponent.Player, neveTheLevelerFillerCreatureUID, neveTheLevelerSetupSrc)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(handBefore)

		deckBefore, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		require.NotEmpty(t, deckBefore)
		wanted := deckBefore[0]

		summonWithOwnMana(t, scn, player, neveTheLevelerUID)

		action, err := scn.LatestAction(player, 0)
		require.NoError(t, err, "expected the search to be offered")
		assert.Equal(t, 0, action.MinSelections, "taking a creature is optional")
		assert.Equal(t, 1, action.MaxSelections, "opponent has exactly one extra creature")

		answerInTurn(t, scn, player, wanted.ID)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount+1)
		assert.Equal(t, match.HAND, wanted.Zone)
	})

	t.Run("several extra opponent creatures allow taking that many", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, opponent.Player, neveTheLevelerFillerCreatureUID, neveTheLevelerSetupSrc)
		putCardInBattlezone(t, scn, opponent.Player, neveTheLevelerFillerCreatureUID, neveTheLevelerSetupSrc)
		putCardInBattlezone(t, scn, opponent.Player, neveTheLevelerFillerCreatureUID, neveTheLevelerSetupSrc)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(handBefore)

		deckBefore, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(deckBefore), 2)
		first, second := deckBefore[0], deckBefore[1]

		summonWithOwnMana(t, scn, player, neveTheLevelerUID)

		action, err := scn.LatestAction(player, 0)
		require.NoError(t, err, "expected the search to be offered")
		assert.Equal(t, 2, action.MaxSelections, "one per extra creature the opponent has")

		answerInTurn(t, scn, player, first.ID, second.ID)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount+2)
		assert.Equal(t, match.HAND, first.Zone)
		assert.Equal(t, match.HAND, second.Zone)
	})

	t.Run("the search may be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, opponent.Player, neveTheLevelerFillerCreatureUID, neveTheLevelerSetupSrc)
		putCardInBattlezone(t, scn, opponent.Player, neveTheLevelerFillerCreatureUID, neveTheLevelerSetupSrc)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(handBefore)

		deckBefore, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		deckCount := len(deckBefore)

		summonWithOwnMana(t, scn, player, neveTheLevelerUID)

		_, err = scn.LatestAction(player, 0)
		require.NoError(t, err, "expected the search to be offered")

		cancelInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount, "nothing was taken")

		deck, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		assert.Len(t, deck, deckCount, "the deck is only shuffled, not shrunk")
	})

	t.Run("a deck without creatures yields nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, opponent.Player, neveTheLevelerFillerCreatureUID, neveTheLevelerSetupSrc)
		putCardInBattlezone(t, scn, opponent.Player, neveTheLevelerFillerCreatureUID, neveTheLevelerSetupSrc)

		player.Player.DestroyDeck()
		for range 10 {
			_, err := player.Player.SpawnCard(neveTheLevelerSpellUID, match.DECK)
			require.NoError(t, err)
		}

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		handCount := len(handBefore)

		summonWithOwnMana(t, scn, player, neveTheLevelerUID)

		// The prompt still opens because a search shows the whole deck, but
		// there is no creature in it to take.
		_, err = scn.LatestAction(player, 0)
		require.NoError(t, err, "expected the search to be offered")

		cancelInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, handCount)
	})
}
