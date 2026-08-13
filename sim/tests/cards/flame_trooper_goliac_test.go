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
	flameTrooperGoliacUID = "34ab516c-6a28-4743-a572-6c1140b1792d"
	flameTrooperSmallUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	flameTrooperHugeUID   = "f6ff8845-afc0-4958-8673-fad12058193a" // Bloodwing Mantis (6000)
	flameTrooperGoliacSrc = "flame_trooper_goliac_test_setup"
)

func TestFlameTrooperGoliac(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, flameTrooperGoliacUID, flameTrooperGoliacSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Flame Trooper Goliac", 4000, 5, []string{civ.Fire})
		assert.True(t, card.HasFamily(family.Armorloid))
		assert.True(t, card.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count it destroys a creature at 5000 or less", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)
		small := putCardInBattlezone(t, scn, opponent.Player, flameTrooperSmallUID, flameTrooperGoliacSrc)
		huge := putCardInBattlezone(t, scn, opponent.Player, flameTrooperHugeUID, flameTrooperGoliacSrc)

		goliac := spawnForLater(t, player, flameTrooperGoliacUID)
		passTurnToSelf(t, scn, player, opponent)
		putIntoPlay(t, scn, player, goliac)

		// Only the small one is in range, so it is taken without asking.
		assert.Equal(t, match.GRAVEYARD, small.Zone)
		assert.Equal(t, match.BATTLEZONE, huge.Zone, "6000 is out of range")
	})

	t.Run("without the count nothing is destroyed", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		small := putCardInBattlezone(t, scn, opponent.Player, flameTrooperSmallUID, flameTrooperGoliacSrc)

		goliac := spawnForLater(t, player, flameTrooperGoliacUID)
		passTurnToSelf(t, scn, player, opponent)
		putIntoPlay(t, scn, player, goliac)

		assert.Equal(t, match.BATTLEZONE, small.Zone)
	})
}
