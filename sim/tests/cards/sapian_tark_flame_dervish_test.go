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
	sapianTarkUID   = "61854910-584d-4d41-a125-61efe3c67a53"
	sapianTarkSetup = "sapian_tark_test_setup"
)

func TestSapianTarkFlameDervish(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		sapian := putCardInBattlezone(t, scn, player.Player, sapianTarkUID, sapianTarkSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, sapian, "Sapian Tark, Flame Dervish", 2000, 3, []string{civ.Fire})
		assert.True(t, sapian.HasFamily(family.Dragonoid))
		assert.True(t, sapian.HasCondition(cnd.WaveStriker))
		assert.False(t, sapian.HasCondition(cnd.AttackUntapped))
	})

	t.Run("with the count it can attack an untapped creature", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		sapian := putCardInBattlezone(t, scn, player.Player, sapianTarkUID, sapianTarkSetup)
		addWaveStrikerFillers(t, scn, player, 2)
		victim := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, sapianTarkSetup)

		passTurnToSelf(t, scn, player, opponent)

		require.True(t, sapian.HasCondition(cnd.AttackUntapped))
		require.False(t, victim.Tapped, "the target is untapped, which normally makes it unattackable")

		require.NoError(t, scn.ActionAttackCreature(player, sapian.ID, victim.ID))

		assert.Equal(t, match.GRAVEYARD, victim.Zone, "6000 beats 2000")
	})
}
