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
	punchTrooperBronksUID = "9f9c1e23-687d-4f8d-a466-5227b45dce40"
	punchSmallUID         = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	punchBigUID           = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	punchSetupSrc         = "punch_trooper_bronks_test_setup"
)

func TestPunchTrooperBronks(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		// Kept in hand, because arriving destroys something.
		card := spawnForLater(t, player, punchTrooperBronksUID)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Punch Trooper Bronks", 3000, 4, []string{civ.Fire})
		assert.True(t, card.HasFamily(family.Armorloid))
	})

	t.Run("the weakest creature on the board is destroyed", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		small := putCardInBattlezone(t, scn, opponent.Player, punchSmallUID, punchSetupSrc)
		big := putCardInBattlezone(t, scn, opponent.Player, punchBigUID, punchSetupSrc)

		bronks := spawnForLater(t, player, punchTrooperBronksUID)
		passTurnToSelf(t, scn, player, opponent)

		// Only one creature is the weakest, so nothing is asked and the move can
		// be made from the test goroutine.
		inPlay := putIntoPlay(t, scn, player, bronks)

		assert.Equal(t, match.GRAVEYARD, small.Zone, "2000 is the lowest power in play")
		assert.Equal(t, match.BATTLEZONE, big.Zone)
		assert.Equal(t, match.BATTLEZONE, inPlay.Zone)
	})

	t.Run("it destroys itself when nothing is weaker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		big := putCardInBattlezone(t, scn, opponent.Player, punchBigUID, punchSetupSrc)

		bronks := spawnForLater(t, player, punchTrooperBronksUID)
		passTurnToSelf(t, scn, player, opponent)

		// Not putIntoPlay: this one is in the graveyard again by the time the
		// move returns, so the battle zone can no longer be asserted on.
		moved, err := player.Player.MoveCard(bronks.ID, match.HAND, match.BATTLEZONE, punchSetupSrc)
		require.NoError(t, err)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, moved.Zone, "3000 is lower than 4000")
		assert.Equal(t, match.BATTLEZONE, big.Zone)
	})

	t.Run("an empty board leaves it alone", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		bronks := spawnForLater(t, player, punchTrooperBronksUID)
		passTurnToSelf(t, scn, player, opponent)

		// It is in the battle zone by the time the effect looks, so it is the
		// only candidate and destroys itself.
		moved, err := player.Player.MoveCard(bronks.ID, match.HAND, match.BATTLEZONE, punchSetupSrc)
		require.NoError(t, err)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, moved.Zone)
	})

	t.Run("a tie is broken by its controller", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirs := putCardInBattlezone(t, scn, opponent.Player, punchSmallUID, punchSetupSrc)
		mine := putCardInBattlezone(t, scn, player.Player, punchSmallUID, punchSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		// Two creatures tied at 2000, so a choice is opened and the summon has
		// to go through the event loop.
		summonWithOwnMana(t, scn, player, punchTrooperBronksUID)

		action, err := scn.WaitForMultipartAction(player, promptStart)
		require.NoError(t, err)
		assert.Contains(t, action.Cards, "Your creatures")
		assert.Contains(t, action.Cards, "Your opponent's creatures")

		answerInTurn(t, scn, player, theirs.ID)
		settleTurn(t, scn)

		assert.Equal(t, match.GRAVEYARD, theirs.Zone)
		assert.Equal(t, match.BATTLEZONE, mine.Zone, "the tie is a choice, not a sweep")
	})
}
