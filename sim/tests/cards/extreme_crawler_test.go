package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	extremeCrawlerUID      = "f98787f3-2962-4f55-80d1-da295cbfede7"
	extremeCrawlerAllyUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (fire, cost 2)
	extremeCrawlerSetupSrc = "extreme_crawler_test_setup"
)

func TestExtremeCrawler(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, extremeCrawlerUID, extremeCrawlerSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Extreme Crawler", 7000, 5, []string{civ.Water})
		assert.True(t, card.HasFamily(family.EarthEater))
		assert.True(t, card.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("arriving clears its controller's other creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		mine := putCardInBattlezone(t, scn, player.Player, extremeCrawlerAllyUID, extremeCrawlerSetupSrc)
		alsoMine := putCardInBattlezone(t, scn, player.Player, extremeCrawlerAllyUID, extremeCrawlerSetupSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, extremeCrawlerAllyUID, extremeCrawlerSetupSrc)

		crawler := spawnForLater(t, player, extremeCrawlerUID)
		passTurnToSelf(t, scn, player, opponent)
		inPlay := putIntoPlay(t, scn, player, crawler)

		assert.Equal(t, match.HAND, mine.Zone)
		assert.Equal(t, match.HAND, alsoMine.Zone)
		assert.Equal(t, match.BATTLEZONE, theirs.Zone, "only its controller's creatures")
		assert.Equal(t, match.BATTLEZONE, inPlay.Zone, "\"other\" spares itself")
	})
}
