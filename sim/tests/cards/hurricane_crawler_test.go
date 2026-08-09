package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	hurricaneCrawlerUID          = "43afea7d-a335-4de4-9eed-a94210d463ce"
	hurricaneCrawlerManaUID      = "9781089f-1aa9-4a75-b106-35e9d431e31d" // Aqua Vehicle
	hurricaneCrawlerFirstHandUID = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	hurricaneCrawlerSecondHandID = "e2b992ee-91a3-49d3-8228-7be60a0b9ec5" // Writhing Bone Ghoul
)

func TestHurricaneCrawler(t *testing.T) {
	t.Run("swaps the whole hand for the same number of mana cards", func(t *testing.T) {
		scn, player, crawler := setupHurricaneCrawlerTest(t)
		player.Player.SpawnCard(hurricaneCrawlerFirstHandUID, match.HAND)
		player.Player.SpawnCard(hurricaneCrawlerSecondHandID, match.HAND)

		firstHandCard, err := scn.FindCard(player.Player, match.HAND, hurricaneCrawlerFirstHandUID)
		require.NoError(t, err)
		secondHandCard, err := scn.FindCard(player.Player, match.HAND, hurricaneCrawlerSecondHandID)
		require.NoError(t, err)

		originalMana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		originalMana = append([]*match.Card(nil), originalMana...)
		require.Len(t, originalMana, 6)

		assert.Equal(t, "Hurricane Crawler", crawler.Name)
		assert.Equal(t, 4000, crawler.Power)
		assert.Equal(t, 5, crawler.ManaCost)
		assert.Equal(t, []string{civ.Water}, crawler.Civs)
		assert.True(t, crawler.HasFamily(family.EarthEater))

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, crawler.ID))
		assert.Equal(t, match.BATTLEZONE, crawler.Zone)

		assert.Equal(t, match.MANAZONE, firstHandCard.Zone)
		assert.Equal(t, match.MANAZONE, secondHandCard.Zone)
		assert.False(t, firstHandCard.Tapped, "Hurricane Crawler does not tap the cards it puts into the mana zone")
		assert.False(t, secondHandCard.Tapped)

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 2, action.MinSelections, "exactly as many cards as were put into the mana zone")
		assert.Equal(t, 2, action.MaxSelections)
		assert.False(t, action.Cancellable)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, originalMana[0].ID, originalMana[1].ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.HAND, originalMana[0].Zone)
		assert.Equal(t, match.HAND, originalMana[1].Zone)
		for _, manaCard := range originalMana[2:] {
			assert.Equal(t, match.MANAZONE, manaCard.Zone)
		}
		assert.Equal(t, match.MANAZONE, firstHandCard.Zone, "cards that were not chosen stay in the mana zone")
		assert.Equal(t, match.MANAZONE, secondHandCard.Zone)

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, 2)
		mana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, mana, 6)
	})

	t.Run("does nothing and opens no prompt with an empty hand", func(t *testing.T) {
		scn, player, crawler := setupHurricaneCrawlerTest(t)

		originalMana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		originalMana = append([]*match.Card(nil), originalMana...)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		// Moved directly so the only card in hand does not have to pay for itself;
		// an effect that prompts here would deadlock the calling goroutine.
		moved, err := player.Player.MoveCard(crawler.ID, match.HAND, match.BATTLEZONE, "hurricane_crawler_test_setup")
		require.NoError(t, err)
		require.Equal(t, match.BATTLEZONE, moved.Zone)

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.NotContains(t, headers, "action")

		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Empty(t, hand)
		for _, manaCard := range originalMana {
			assert.Equal(t, match.MANAZONE, manaCard.Zone)
		}
	})

	t.Run("returns only what remains when a nested effect shrinks the mana zone", func(t *testing.T) {
		scn := scenario.New()
		player := scn.Match.CurrentPlayer()
		clearHurricaneCrawlerTestZone(t, player.Player, match.HAND)

		player.Player.SpawnCard(hurricaneCrawlerUID, match.HAND)
		for range 5 {
			player.Player.SpawnCard(hurricaneCrawlerManaUID, match.MANAZONE)
		}
		for range 3 {
			player.Player.SpawnCard(hurricaneCrawlerFirstHandUID, match.HAND)
		}

		crawler, err := scn.FindCard(player.Player, match.HAND, hurricaneCrawlerUID)
		require.NoError(t, err)

		mana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		require.Len(t, mana, 5)

		// Every MoveCard into the mana zone dispatches events another card can
		// react to. This saboteur empties the mana zone in the middle of the
		// swap, so the card ends up owing more cards than the zone can supply.
		arrivals := 0
		fired := false
		mana[0].Use(func(card *match.Card, ctx *match.Context) {
			if fired {
				return
			}

			event, ok := ctx.Event.(*match.CardMoved)
			if !ok || event.To != match.MANAZONE {
				return
			}

			arrivals++
			if arrivals < 2 {
				return
			}
			fired = true

			current, err := card.Player.Container(match.MANAZONE)
			if err != nil {
				return
			}

			for _, manaCard := range append([]*match.Card(nil), current...) {
				if manaCard.ID == event.CardID {
					continue
				}

				card.Player.MoveCard(manaCard.ID, match.MANAZONE, match.GRAVEYARD, "hurricane_crawler_test_saboteur")
			}
		})

		require.NoError(t, scn.ActionPlayCard(player, crawler.ID))
		assert.Equal(t, match.BATTLEZONE, crawler.Zone)

		// Three cards left the hand but only two survived in the mana zone, so
		// the selection resolves with those two instead of stalling on a prompt
		// that can never be satisfied.
		remainingMana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Empty(t, remainingMana)
		hand, err := player.Player.Container(match.HAND)
		require.NoError(t, err)
		assert.Len(t, hand, 2)

		// The sequential event loop is still healthy after the clamped selection.
		require.NoError(t, scn.ActionEndTurn(player))
	})

	t.Run("puts a single card back when only one card is swapped", func(t *testing.T) {
		scn, player, crawler := setupHurricaneCrawlerTest(t)
		player.Player.SpawnCard(hurricaneCrawlerFirstHandUID, match.HAND)

		handCard, err := scn.FindCard(player.Player, match.HAND, hurricaneCrawlerFirstHandUID)
		require.NoError(t, err)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(player, crawler.ID))

		action, err := scn.LatestAction(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, action.MinSelections)
		assert.Equal(t, 1, action.MaxSelections)

		completionStart, err := scn.MessageCount(player)
		require.NoError(t, err)
		require.NoError(t, scn.SubmitAction(player, handCard.ID))
		require.NoError(t, scn.WaitForMessage(player, completionStart, "state_update"))

		assert.Equal(t, match.HAND, handCard.Zone, "the card just put into mana may be chosen again")
		mana, err := player.Player.Container(match.MANAZONE)
		require.NoError(t, err)
		assert.Len(t, mana, 6)
	})
}

func setupHurricaneCrawlerTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.Card) {
	t.Helper()

	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	clearHurricaneCrawlerTestZone(t, player.Player, match.HAND)

	player.Player.SpawnCard(hurricaneCrawlerUID, match.HAND)
	for range 6 {
		player.Player.SpawnCard(hurricaneCrawlerManaUID, match.MANAZONE)
	}

	crawler, err := scn.FindCard(player.Player, match.HAND, hurricaneCrawlerUID)
	require.NoError(t, err)

	return scn, player, crawler
}

func clearHurricaneCrawlerTestZone(t *testing.T, player *match.Player, zone string) {
	t.Helper()

	cards, err := player.Container(zone)
	require.NoError(t, err)
	for _, card := range append([]*match.Card(nil), cards...) {
		moved, err := player.MoveCard(card.ID, zone, match.GRAVEYARD, "hurricane_crawler_test_setup")
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)
	}
}
