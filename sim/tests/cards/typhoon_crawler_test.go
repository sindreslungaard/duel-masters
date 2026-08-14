package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	typhoonCrawlerUID       = "aeaaf98d-938f-46d1-a271-49a86f668ae6"
	typhoonFireAttackerUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (fire, cost 2)
	typhoonWaterAttackerUID = "f4a364f5-d0e9-4777-b51e-6dc6e39b803c" // Aqua Shooter (water)
	typhoonCrawlerSetupSrc  = "typhoon_crawler_test_setup"
)

func TestTyphoonCrawler(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, typhoonCrawlerUID, typhoonCrawlerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Typhoon Crawler", 5000, 6, []string{civ.Water})
		assert.True(t, card.HasFamily(family.EarthEater))
	})

	t.Run("fire creatures cannot attack it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		crawler := putCardInBattlezone(t, scn, player.Player, typhoonCrawlerUID, typhoonCrawlerSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, typhoonFireAttackerUID, typhoonCrawlerSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		crawler.Tapped = true

		require.Error(t, scn.ActionAttackCreature(opponent, attacker.ID, crawler.ID), "a fire attacker is refused")
		assert.Equal(t, match.BATTLEZONE, crawler.Zone)
	})

	t.Run("other civilizations can attack it", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		crawler := putCardInBattlezone(t, scn, player.Player, typhoonCrawlerUID, typhoonCrawlerSetupSrc)
		attacker := putCardInBattlezone(t, scn, opponent.Player, typhoonWaterAttackerUID, typhoonCrawlerSetupSrc)

		require.NoError(t, scn.ActionEndTurn(player))
		crawler.Tapped = true

		require.NoError(t, scn.ActionAttackCreature(opponent, attacker.ID, crawler.ID), "water is not covered")
	})
}
