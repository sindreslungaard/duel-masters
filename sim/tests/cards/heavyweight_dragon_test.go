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
	heavyweightDragonUID = "c2970e00-4951-421c-837a-119a3bf564d8"
	heavyweightSmallUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg (2000)
	heavyweightFourUID   = "8112be9d-50a9-4489-b3f8-257aeed62205" // Magmadragon Melgars (4000)
	heavyweightFiveUID   = "84e1b416-c2d5-4ae1-aca0-025651c6aa58" // Tri-horn Shepherd (5000)
	heavyweightSetupSrc  = "heavyweight_dragon_test_setup"
)

func TestHeavyweightDragon(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dragon := putCardInBattlezone(t, scn, player.Player, heavyweightDragonUID, heavyweightSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, dragon, "Heavyweight Dragon", 9000, 7, []string{civ.Fire})
		assert.True(t, dragon.HasFamily(family.ArmoredDragon))
		assert.True(t, dragon.HasCondition(cnd.DoubleBreaker))
		assert.True(t, dragon.HasCondition(cnd.TapAbility))
	})

	t.Run("two tapped creatures under its power are both destroyed", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dragon := putCardInBattlezone(t, scn, player.Player, heavyweightDragonUID, heavyweightSetupSrc)
		first := putCardInBattlezone(t, scn, opponent.Player, heavyweightSmallUID, heavyweightSetupSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, heavyweightSmallUID, heavyweightSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		first.Tapped = true
		second.Tapped = true

		require.NoError(t, scn.ActionUseTapAbility(player, dragon.ID))
		answerInTurn(t, scn, player, first.ID, second.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, first.Zone, "4000 total is under 9000")
		assert.Equal(t, match.GRAVEYARD, second.Zone)
		assert.True(t, dragon.Tapped, "using the tap ability taps it")
	})

	t.Run("a selection that reaches its power is explained and offered again", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dragon := putCardInBattlezone(t, scn, player.Player, heavyweightDragonUID, heavyweightSetupSrc)
		first := putCardInBattlezone(t, scn, opponent.Player, heavyweightFourUID, heavyweightSetupSrc)
		second := putCardInBattlezone(t, scn, opponent.Player, heavyweightFiveUID, heavyweightSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		first.Tapped = true
		second.Tapped = true
		require.Equal(t, 9000, scn.Match.GetPower(first, false)+scn.Match.GetPower(second, false))

		warningStart, err := scn.MessageCount(player)
		require.NoError(t, err)

		require.NoError(t, scn.ActionUseTapAbility(player, dragon.ID))
		answerInTurn(t, scn, player, first.ID, second.ID)

		// The combined power has to be strictly less, so an exact match fails.
		warnings, err := scn.Warnings(player, warningStart)
		require.NoError(t, err)
		require.NotEmpty(t, warnings, "the player should be told why it did not work")
		assert.Contains(t, warnings[0], "9000 total power", "the warning shows the sum")
		assert.Contains(t, warnings[0], "Magmadragon Melgars (4000)", "and how each card contributed")

		// The selection is open again; giving up leaves the board alone.
		cancelInTurn(t, scn, player)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, first.Zone)
		assert.Equal(t, match.BATTLEZONE, second.Zone)
		assert.True(t, dragon.Tapped, "the ability was still used")
	})

	t.Run("a second attempt after a warning can still succeed", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dragon := putCardInBattlezone(t, scn, player.Player, heavyweightDragonUID, heavyweightSetupSrc)
		heavy := putCardInBattlezone(t, scn, opponent.Player, heavyweightFourUID, heavyweightSetupSrc)
		heavier := putCardInBattlezone(t, scn, opponent.Player, heavyweightFiveUID, heavyweightSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		heavy.Tapped = true
		heavier.Tapped = true

		require.NoError(t, scn.ActionUseTapAbility(player, dragon.ID))

		// Too heavy together, but one of them on its own is well under.
		answerInTurn(t, scn, player, heavy.ID, heavier.ID)
		answerInTurn(t, scn, player, heavy.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, heavy.Zone)
		assert.Equal(t, match.BATTLEZONE, heavier.Zone)
	})

	t.Run("untapped creatures are not offered", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dragon := putCardInBattlezone(t, scn, player.Player, heavyweightDragonUID, heavyweightSetupSrc)
		untapped := putCardInBattlezone(t, scn, opponent.Player, heavyweightSmallUID, heavyweightSetupSrc)
		tapped := putCardInBattlezone(t, scn, opponent.Player, heavyweightSmallUID, heavyweightSetupSrc)

		passTurnToSelf(t, scn, player, opponent)
		tapped.Tapped = true

		require.NoError(t, scn.ActionUseTapAbility(player, dragon.ID))

		// "Up to 2" is declinable, so even a single candidate is offered rather
		// than taken automatically. The untapped creature is not among them.
		action, err := scn.LatestAction(player, 0)
		require.NoError(t, err)
		require.Len(t, action.Cards, 1, "only the tapped creature is offered")

		answerInTurn(t, scn, player, tapped.ID)
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.GRAVEYARD, tapped.Zone)
		assert.Equal(t, match.BATTLEZONE, untapped.Zone)
	})

	t.Run("nothing tapped means nothing to choose", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		dragon := putCardInBattlezone(t, scn, player.Player, heavyweightDragonUID, heavyweightSetupSrc)
		untapped := putCardInBattlezone(t, scn, opponent.Player, heavyweightSmallUID, heavyweightSetupSrc)

		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionUseTapAbility(player, dragon.ID))
		require.NoError(t, scn.WaitForEventLoop())

		assert.Equal(t, match.BATTLEZONE, untapped.Zone)
		assert.True(t, dragon.Tapped)
	})
}
