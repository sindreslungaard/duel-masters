package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"testing"

	"github.com/stretchr/testify/assert"
)

const bradsCutterUID = "49373744-6cbe-4247-8ae4-12fcc0f62201"

func TestBradsCutter(t *testing.T) {
	scn, player, opponent := setupDuel(t)
	cutter := putCardInBattlezone(t, scn, player.Player, bradsCutterUID, "brads_cutter_test_setup")
	passTurnToSelf(t, scn, player, opponent)

	assertPrinted(t, cutter, "Brad's Cutter", 1000, 2, []string{civ.Fire})
	assert.True(t, cutter.HasFamily(family.Xenoparts))
	assert.True(t, cutter.HasCondition(cnd.ShieldTrigger))
}
