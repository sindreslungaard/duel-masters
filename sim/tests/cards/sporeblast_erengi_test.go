package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	sporeblastErengiUID         = "19e27fb6-f221-44ae-9e22-ec48829ac117"
	sporeblastErengiCreatureUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	sporeblastErengiSpellUID    = "5883180e-d88c-4f24-b17c-f5a837420147" // Terror Pit
	sporeblastErengiSetupSrc    = "sporeblast_erengi_test_setup"
)

func TestSporeblastErengi(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		erengi := putCardInBattlezone(t, scn, player.Player, sporeblastErengiUID, sporeblastErengiSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Sporeblast Erengi", erengi.Name)
		assert.Equal(t, 4000, erengi.Power)
		assert.Equal(t, 4, erengi.ManaCost)
		assert.Equal(t, []string{civ.Nature}, erengi.Civs)
		assert.Equal(t, []string{civ.Nature}, erengi.ManaRequirement)
		assert.True(t, erengi.HasFamily(family.BalloonMushroom))
		assert.True(t, erengi.HasCondition(cnd.SilentSkill))
	})

	t.Run("takes a creature out of the deck", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		erengi := putCardInBattlezone(t, scn, player.Player, sporeblastErengiUID, sporeblastErengiSetupSrc)
		erengi.Tapped = true

		passTurnToSelf(t, scn, player, opponent)

		deckBefore, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		require.NotEmpty(t, deckBefore)

		wanted := deckBefore[0]

		useSilentSkill(t, scn, player)
		require.NoError(t, scn.SubmitAction(player, wanted.ID))
		settleTurn(t, scn)

		assert.Equal(t, match.HAND, wanted.Zone)

		deck, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		// Silent skill resolves at the start of the turn, so the draw step
		// still follows: one card for the search and one for the draw.
		assert.Len(t, deck, len(deckBefore)-2)
		assert.Len(t, hand, len(handBefore)+2)
		assert.True(t, erengi.Tapped)
	})

	t.Run("the search may be declined", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		erengi := putCardInBattlezone(t, scn, player.Player, sporeblastErengiUID, sporeblastErengiSetupSrc)
		erengi.Tapped = true

		passTurnToSelf(t, scn, player, opponent)

		deckBefore, err := player.Player.Container(match.DECK)
		require.NoError(t, err)

		useSilentSkill(t, scn, player)
		require.NoError(t, scn.CancelAction(player))
		settleTurn(t, scn)

		deck, err := player.Player.Container(match.DECK)
		require.NoError(t, err)
		assert.Len(t, deck, len(deckBefore)-1, "only the draw step, since the printed wording is \"you may take\"")
		assert.True(t, erengi.Tapped)
	})

	t.Run("a deck without creatures yields nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		erengi := putCardInBattlezone(t, scn, player.Player, sporeblastErengiUID, sporeblastErengiSetupSrc)
		erengi.Tapped = true

		player.Player.DestroyDeck()
		for range 4 {
			_, err := player.Player.SpawnCard(sporeblastErengiSpellUID, match.DECK)
			require.NoError(t, err)
		}

		passTurnToSelf(t, scn, player, opponent)

		handBefore, err := player.Player.Container(match.HAND)
		require.NoError(t, err)

		useSilentSkill(t, scn, player)

		// The prompt still opens, because a search shows the whole deck, but
		// there is no creature in it to take.
		require.NoError(t, scn.CancelAction(player))
		settleTurn(t, scn)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, len(handBefore)+1, "only the draw step")
	})

	t.Run("it cannot search the opponent's deck", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		erengi := putCardInBattlezone(t, scn, player.Player, sporeblastErengiUID, sporeblastErengiSetupSrc)
		erengi.Tapped = true

		passTurnToSelf(t, scn, player, opponent)

		opponentDeckBefore, err := opponent.Player.Container(match.DECK)
		require.NoError(t, err)

		useSilentSkill(t, scn, player)
		require.NoError(t, scn.CancelAction(player))
		settleTurn(t, scn)

		opponentDeck, err := opponent.Player.Container(match.DECK)
		require.NoError(t, err)
		assert.Len(t, opponentDeck, len(opponentDeckBefore))
	})
}
