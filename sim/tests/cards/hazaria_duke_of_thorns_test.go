package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	hazariaUID   = "73af7f15-12fc-4ae7-a7c3-08d7141ad818"
	hazariaSetup = "hazaria_duke_of_thorns_test_setup"
)

func TestHazariaDukeOfThorns(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		hazaria := putCardInBattlezone(t, scn, player.Player, hazariaUID, hazariaSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, hazaria, "Hazaria, Duke of Thorns", 2000, 4, []string{civ.Darkness})
		assert.True(t, hazaria.HasFamily(family.DarkLord))
		assert.True(t, hazaria.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count the opponent loses a creature of their choosing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)
		theirs := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, hazariaSetup)

		hazaria := spawnForLater(t, player, hazariaUID)
		passTurnToSelf(t, scn, player, opponent)
		putIntoPlay(t, scn, player, hazaria)

		// One creature, so the choice resolves without a prompt.
		assert.Equal(t, match.GRAVEYARD, theirs.Zone)
	})

	t.Run("without the count nothing happens", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirs := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, hazariaSetup)

		hazaria := spawnForLater(t, player, hazariaUID)
		passTurnToSelf(t, scn, player, opponent)
		putIntoPlay(t, scn, player, hazaria)

		assert.Equal(t, match.BATTLEZONE, theirs.Zone)
	})
}
