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
	clonedSpikeHornUID = "3bd89f19-36f5-4f39-ab3c-4ecc3a4204a0"
	clonedSpikeHornSrc = "cloned_spike_horn_test_setup"
)

func TestClonedSpikeHorn(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, clonedSpikeHornUID, clonedSpikeHornSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Cloned Spike-Horn", 3000, 4, []string{civ.Nature})
		assert.True(t, card.HasFamily(family.HornedBeast))
	})

	t.Run("it grows by every copy in either graveyard", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		spikeHorn := putCardInBattlezone(t, scn, player.Player, clonedSpikeHornUID, clonedSpikeHornSrc)

		assert.Equal(t, 3000, scn.Match.GetPower(spikeHorn, false), "it never counts itself")

		_, err := player.Player.SpawnCard(clonedSpikeHornUID, match.GRAVEYARD)
		require.NoError(t, err)
		assert.Equal(t, 6000, scn.Match.GetPower(spikeHorn, false))

		_, err = opponent.Player.SpawnCard(clonedSpikeHornUID, match.GRAVEYARD)
		require.NoError(t, err)
		assert.Equal(t, 9000, scn.Match.GetPower(spikeHorn, false), "each graveyard means each")
	})

	t.Run("copies anywhere else are worth nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		spikeHorn := putCardInBattlezone(t, scn, player.Player, clonedSpikeHornUID, clonedSpikeHornSrc)

		for _, zone := range []string{match.HAND, match.MANAZONE, match.SHIELDZONE} {
			_, err := player.Player.SpawnCard(clonedSpikeHornUID, zone)
			require.NoError(t, err)
		}

		putCardInBattlezone(t, scn, player.Player, clonedSpikeHornUID, clonedSpikeHornSrc)

		_, err := opponent.Player.SpawnCard(clonedSpikeHornUID, match.HAND)
		require.NoError(t, err)

		assert.Equal(t, 3000, scn.Match.GetPower(spikeHorn, false))
	})

	t.Run("a single copy is enough for double breaker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		spikeHorn := putCardInBattlezone(t, scn, player.Player, clonedSpikeHornUID, clonedSpikeHornSrc)

		_, err := player.Player.SpawnCard(clonedSpikeHornUID, match.GRAVEYARD)
		require.NoError(t, err)

		passTurnToSelf(t, scn, player, opponent)

		require.Equal(t, 6000, scn.Match.GetPower(spikeHorn, false))
		assert.True(t, spikeHorn.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("alone it stays a single breaker", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		spikeHorn := putCardInBattlezone(t, scn, player.Player, clonedSpikeHornUID, clonedSpikeHornSrc)

		passTurnToSelf(t, scn, player, opponent)

		assert.False(t, spikeHorn.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("a copy in the graveyard is not lifted by the ones beside it", func(t *testing.T) {
		scn, player, _ := setupDuel(t)

		buried, err := player.Player.SpawnCard(clonedSpikeHornUID, match.GRAVEYARD)
		require.NoError(t, err)
		_, err = player.Player.SpawnCard(clonedSpikeHornUID, match.GRAVEYARD)
		require.NoError(t, err)

		// Two in the graveyard, and the one being asked about only sees the
		// other, never itself.
		assert.Equal(t, 6000, scn.Match.GetPower(buried, false))
	})
}
