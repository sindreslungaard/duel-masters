package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fantasyFishUID = "af369960-c65f-49f8-bb21-aa6190d0b3b4"

func TestFantasyFish(t *testing.T) {
	scn, player, opponent := setupDuel(t)
	fish := putCardInBattlezone(t, scn, player.Player, fantasyFishUID, "fantasy_fish_test_setup")
	passTurnToSelf(t, scn, player, opponent)

	assertPrinted(t, fish, "Fantasy Fish", 2000, 7, []string{civ.Water})
	assert.True(t, fish.HasFamily(family.GelFish))
	assert.True(t, fish.HasCondition(cnd.ShieldTrigger))
	assert.True(t, fish.HasCondition(cnd.Blocker))
}
