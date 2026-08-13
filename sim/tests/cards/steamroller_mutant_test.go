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
	steamrollerMutantUID      = "6250a4d5-c34e-4585-9c3c-cda7b0b094d3"
	steamrollerVictimUID      = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	steamrollerMutantSetupSrc = "steamroller_mutant_test_setup"
)

func TestSteamrollerMutant(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, steamrollerMutantUID, steamrollerMutantSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Steamroller Mutant", 3000, 4, []string{civ.Darkness})
		assert.True(t, card.HasFamily(family.Hedrian))
		assert.True(t, card.HasCondition(cnd.WaveStriker))
	})

	t.Run("with the count it clears the whole board", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		addWaveStrikerFillers(t, scn, player, 2)
		theirs := putCardInBattlezone(t, scn, opponent.Player, steamrollerVictimUID, steamrollerMutantSetupSrc)

		mutant := spawnForLater(t, player, steamrollerMutantUID)
		passTurnToSelf(t, scn, player, opponent)

		moved, err := player.Player.MoveCard(mutant.ID, match.HAND, match.BATTLEZONE, steamrollerMutantSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, theirs.Zone)
		assert.Equal(t, match.GRAVEYARD, moved.Zone, "it destroys itself too")

		mine, err := player.Player.Container(match.BATTLEZONE)
		require.NoError(t, err)
		assert.Empty(t, mine, "including the wave strikers that switched it on")
	})

	t.Run("without the count nothing is destroyed", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirs := putCardInBattlezone(t, scn, opponent.Player, steamrollerVictimUID, steamrollerMutantSetupSrc)

		mutant := spawnForLater(t, player, steamrollerMutantUID)
		passTurnToSelf(t, scn, player, opponent)
		inPlay := putIntoPlay(t, scn, player, mutant)

		assert.Equal(t, match.BATTLEZONE, theirs.Zone)
		assert.Equal(t, match.BATTLEZONE, inPlay.Zone)
	})
}
