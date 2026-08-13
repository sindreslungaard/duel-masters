package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	mizoyTheOracleUID      = "899e15d3-ed0d-4f13-adf9-54e2db77736d"
	mizoyFireCreatureUID   = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (fire, cost 2)
	mizoyLightCreatureUID  = "d0b9d738-bd33-4dd1-b6dc-234897f52266" // Belmol, the Explorer (light)
	mizoyTheOracleSetupSrc = "mizoy_the_oracle_test_setup"
)

func TestMizoyTheOracle(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		// Kept in hand: it may ask for a target as it arrives, and a card moved into play from the test
		// goroutine would leave it waiting on a prompt only it could answer.
		card := spawnForLater(t, player, mizoyTheOracleUID)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Mizoy, the Oracle", 2500, 3, []string{civ.Light})
		assert.True(t, card.HasFamily(family.LightBringer))
	})

	t.Run("it may tap a darkness or fire creature on either side", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirFire := putCardInBattlezone(t, scn, opponent.Player, mizoyFireCreatureUID, mizoyTheOracleSetupSrc)

		summonWithOwnMana(t, scn, player, mizoyTheOracleUID)

		answerInTurn(t, scn, player, theirFire.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.True(t, theirFire.Tapped)
	})

	t.Run("nothing matching means no prompt", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		theirLight := putCardInBattlezone(t, scn, opponent.Player, mizoyLightCreatureUID, mizoyTheOracleSetupSrc)

		promptStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		summonWithOwnMana(t, scn, player, mizoyTheOracleUID)
		require.NoError(t, scn.WaitForEventLoop())

		headers, err := scn.MessageHeaders(player, promptStart)
		require.NoError(t, err)
		assert.Equal(t, 1, countHeaders(headers, "action"), "only the mana payment should ask anything")
		assert.False(t, theirLight.Tapped)
	})
}
