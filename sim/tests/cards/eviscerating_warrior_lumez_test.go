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
	lumezUID   = "ad3e67c6-5f09-4602-8586-53cc80813555"
	lumezSetup = "eviscerating_warrior_lumez_test_setup"
)

func TestEvisceratingWarriorLumez(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		lumez := putCardInBattlezone(t, scn, player.Player, lumezUID, lumezSetup)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, lumez, "Eviscerating Warrior Lumez", 2000, 3, []string{civ.Fire})
		assert.True(t, lumez.HasFamily(family.Armorloid))
		assert.True(t, lumez.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count it sweeps everything at 2000 or less, itself included", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)
		theirSmall := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, lumezSetup)
		theirBig := putCardInBattlezone(t, scn, opponent.Player, waveStrikerBigCreatureUID, lumezSetup)

		lumez := spawnForLater(t, player, lumezUID)
		passTurnToSelf(t, scn, player, opponent)

		// Moved by hand rather than through putIntoPlay: the sweep runs while
		// the move is still resolving, so it is already in the graveyard by the
		// time MoveCard returns.
		_, err := player.Player.MoveCard(lumez.ID, match.HAND, match.BATTLEZONE, lumezSetup)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())
		inPlay := lumez

		assert.Equal(t, match.GRAVEYARD, theirSmall.Zone)
		assert.Equal(t, match.BATTLEZONE, theirBig.Zone)
		assert.Equal(t, match.GRAVEYARD, inPlay.Zone, "at 2000 power it is caught by its own sweep")
	})

	t.Run("without the count nothing is destroyed", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirSmall := putCardInBattlezone(t, scn, opponent.Player, immortalBaronVorgUID, lumezSetup)

		lumez := spawnForLater(t, player, lumezUID)
		passTurnToSelf(t, scn, player, opponent)
		inPlay := putIntoPlay(t, scn, player, lumez)

		assert.Equal(t, match.BATTLEZONE, theirSmall.Zone)
		assert.Equal(t, match.BATTLEZONE, inPlay.Zone)
	})
}
