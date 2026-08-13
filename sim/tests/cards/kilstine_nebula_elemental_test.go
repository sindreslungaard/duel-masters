package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	kilstineUID      = "af2e10ca-da3c-48ee-8064-d12c400ff1f9"
	kilstineAllyUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	kilstineSetupSrc = "kilstine_nebula_elemental_test_setup"
)

func TestKilstineNebulaElemental(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, kilstineUID, kilstineSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Kilstine, Nebula Elemental", 5000, 7, []string{civ.Light})
		assert.True(t, card.HasFamily(family.AngelCommand))
		assert.True(t, card.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count it lifts its other creatures", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		kilstine := putCardInBattlezone(t, scn, player.Player, kilstineUID, kilstineSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, kilstineAllyUID, kilstineSetupSrc)
		theirs := putCardInBattlezone(t, scn, opponent.Player, kilstineAllyUID, kilstineSetupSrc)
		addWaveStrikerFillers(t, scn, player, 2)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 7000, scn.Match.GetPower(ally, false), "2000 plus 5000")
		assert.True(t, ally.HasCondition(cnd.Blocker))
		assert.True(t, ally.HasCondition(cnd.DoubleBreaker))

		assert.Equal(t, 5000, scn.Match.GetPower(kilstine, false), "\"other\" spares itself")
		assert.False(t, kilstine.HasCondition(cnd.DoubleBreaker))

		assert.Equal(t, 2000, scn.Match.GetPower(theirs, false), "only its controller's creatures")
	})

	t.Run("without the count nothing is granted", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, kilstineUID, kilstineSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, kilstineAllyUID, kilstineSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, 2000, scn.Match.GetPower(ally, false))
		assert.False(t, ally.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("the grant is withdrawn when the count drops", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		putCardInBattlezone(t, scn, player.Player, kilstineUID, kilstineSetupSrc)
		ally := putCardInBattlezone(t, scn, player.Player, kilstineAllyUID, kilstineSetupSrc)
		filler := putCardInBattlezone(t, scn, player.Player, waveStrikerFillerUID, kilstineSetupSrc)
		putCardInBattlezone(t, scn, player.Player, waveStrikerFillerUID, kilstineSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		require.Equal(t, 7000, scn.Match.GetPower(ally, false))

		_, err := player.Player.MoveCard(filler.ID, match.BATTLEZONE, match.GRAVEYARD, kilstineSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, 2000, scn.Match.GetPower(ally, false))
		assert.False(t, ally.HasCondition(cnd.DoubleBreaker))
	})
}
