package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	heartyCapnPolligonUID = "2967d4cc-b4da-4511-b9e7-95b4844216f9"
	polligonSetupSrc      = "hearty_capn_polligon_test_setup"
)

func TestHeartyCapnPolligon(t *testing.T) {
	t.Run("printed characteristics", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		polligon := putCardInBattlezone(t, scn, player.Player, heartyCapnPolligonUID, polligonSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		assertPrinted(t, polligon, "Hearty Cap'n Polligon", 2000, 1, []string{civ.Nature})
		assert.True(t, polligon.HasFamily(family.SnowFaerie))
	})

	t.Run("it comes home after breaking a shield", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		polligon := putCardInBattlezone(t, scn, player.Player, heartyCapnPolligonUID, polligonSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		breakAShieldWith(t, scn, player, opponent, polligon)
		require.Equal(t, match.BATTLEZONE, polligon.Zone, "it stays out until the end of the turn")

		require.NoError(t, scn.ActionEndTurn(player))

		assert.Equal(t, match.HAND, polligon.Zone)
	})

	t.Run("it stays out on a turn it broke nothing", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		polligon := putCardInBattlezone(t, scn, player.Player, heartyCapnPolligonUID, polligonSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		require.NoError(t, scn.ActionEndTurn(player))

		assert.Equal(t, match.BATTLEZONE, polligon.Zone)
	})

	t.Run("the memory of breaking does not carry into a later turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		polligon := putCardInBattlezone(t, scn, player.Player, heartyCapnPolligonUID, polligonSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		breakAShieldWith(t, scn, player, opponent, polligon)
		require.NoError(t, scn.ActionEndTurn(player))
		require.Equal(t, match.HAND, polligon.Zone)

		// Put back out without attacking: the flag from the earlier turn must
		// not bounce it a second time.
		require.NoError(t, scn.ActionEndTurn(opponent))
		returned, err := player.Player.MoveCard(polligon.ID, match.HAND, match.BATTLEZONE, polligonSetupSrc)
		require.NoError(t, err)
		require.NoError(t, scn.WaitForEventLoop())

		require.NoError(t, scn.ActionEndTurn(player))

		assert.Equal(t, match.BATTLEZONE, returned.Zone)
	})

	t.Run("it is not returned at the end of the opponent's turn", func(t *testing.T) {
		scn, player, opponent := setupDuel(t)
		polligon := putCardInBattlezone(t, scn, player.Player, heartyCapnPolligonUID, polligonSetupSrc)
		passTurnToSelf(t, scn, player, opponent)

		breakAShieldWith(t, scn, player, opponent, polligon)
		require.NoError(t, scn.ActionEndTurn(player))
		require.Equal(t, match.HAND, polligon.Zone)

		// The opponent's own turn ending must not do anything to it.
		require.NoError(t, scn.ActionEndTurn(opponent))
		assert.Equal(t, match.HAND, polligon.Zone)
	})
}

// breakAShieldWith attacks the opponent directly and breaks one shield.
func breakAShieldWith(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, opponent *match.PlayerReference, attacker *match.Card) {
	t.Helper()

	shields, err := opponent.Player.Container(match.SHIELDZONE)
	require.NoError(t, err)
	require.NotEmpty(t, shields)

	_, err = scn.ActionAttackPlayer(player, attacker.ID)
	require.NoError(t, err)
	require.NoError(t, scn.ResolveAttack(player, shields[0].ID))
	require.NoError(t, scn.WaitForEventLoop())

	require.Equal(t, match.HAND, shields[0].Zone, "the shield really was broken")
}
