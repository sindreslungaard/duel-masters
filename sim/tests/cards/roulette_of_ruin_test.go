package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	rouletteOfRuinUID = "c447753f-ca24-4482-8db0-837dd6f7d31b"
	rouletteCost2UID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (cost 2)
	rouletteCost4UID  = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (cost 4)
	rouletteSetupSrc  = "roulette_of_ruin_test_setup"
)

func TestRouletteOfRuin(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		spell, err := player.Player.SpawnCard(rouletteOfRuinUID, match.HAND)
		require.NoError(t, err)
		passTurnToSelf(t, scn, player, opponent)

		assert.Equal(t, "Roulette of Ruin", spell.Name)
		assert.Equal(t, 5, spell.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, spell.Civs)
		assert.True(t, spell.HasCondition(cnd.ShieldTrigger))
	})

	t.Run("both players lose every card of the chosen cost", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		emptyHand(t, player, rouletteSetupSrc)
		emptyHand(t, opponent, rouletteSetupSrc)

		myTwo, err := player.Player.SpawnCard(rouletteCost2UID, match.HAND)
		require.NoError(t, err)
		myFour, err := player.Player.SpawnCard(rouletteCost4UID, match.HAND)
		require.NoError(t, err)
		theirTwo, err := opponent.Player.SpawnCard(rouletteCost2UID, match.HAND)
		require.NoError(t, err)
		theirFour, err := opponent.Player.SpawnCard(rouletteCost4UID, match.HAND)
		require.NoError(t, err)

		castSpell(t, scn, player, rouletteOfRuinUID)

		// The number is chosen before either hand is shown.
		require.NoError(t, scn.SubmitCount(player, 2))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, myTwo.Zone, "its caster is not spared")
		assert.Equal(t, match.GRAVEYARD, theirTwo.Zone)
		assert.Equal(t, match.HAND, myFour.Zone)
		assert.Equal(t, match.HAND, theirFour.Zone)
	})

	t.Run("a number nobody holds costs nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)

		emptyHand(t, player, rouletteSetupSrc)
		emptyHand(t, opponent, rouletteSetupSrc)

		mine, err := player.Player.SpawnCard(rouletteCost2UID, match.HAND)
		require.NoError(t, err)
		theirs, err := opponent.Player.SpawnCard(rouletteCost2UID, match.HAND)
		require.NoError(t, err)

		castSpell(t, scn, player, rouletteOfRuinUID)
		require.NoError(t, scn.SubmitCount(player, 7))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.HAND, mine.Zone)
		assert.Equal(t, match.HAND, theirs.Zone)
	})
}
