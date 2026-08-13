package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"testing"

	"github.com/stretchr/testify/assert"
)

const melodicHunterUID = "6d65304b-30ab-4943-a88d-dbcd8a204c2c"

func TestMelodicHunter(t *testing.T) {
	scn, player, opponent := setupDuel(t)
	hunter := putCardInBattlezone(t, scn, player.Player, melodicHunterUID, "melodic_hunter_test_setup")
	passTurnToSelf(t, scn, player, opponent)

	assertPrinted(t, hunter, "Melodic Hunter", 3000, 5, []string{civ.Water})
	assert.True(t, hunter.HasFamily(family.Merfolk))
	assert.True(t, hunter.HasCondition(cnd.Blocker))
}
