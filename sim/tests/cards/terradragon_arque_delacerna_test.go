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
	terradragonArqueUID      = "dd02d138-9915-4204-856f-d20427072339"
	terradragonArqueSetupSrc = "terradragon_arque_delacerna_test_setup"
)

func TestTerradragonArqueDelacerna(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		card := putCardInBattlezone(t, scn, player.Player, terradragonArqueUID, terradragonArqueSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, card, "Terradragon Arque Delacerna", 6000, 8, []string{civ.Nature})
		assert.True(t, card.HasFamily(family.EarthDragon))
		assert.True(t, card.HasCondition(cnd.DoubleBreaker))
	})

	t.Run("a discard on the opponent's turn may put it into play instead", func(t *testing.T) {
		scn, player, opponent, dragon := setupDiscardReplacementTest(t, terradragonArqueUID)

		require.NoError(t, scn.ActionPlayCard(opponent, discardTriggerCard(t, scn, opponent)))
		require.NoError(t, scn.SubmitAction(player))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, dragon.Zone, "an 8 cost creature arrives for free")
	})
}
