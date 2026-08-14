package cards

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	gabzagulWarlordOfPainUID      = "d53bb7b9-5d3e-44ed-a8ab-ab262ec23cb8"
	gabzagulWarlordOfPainAllyUID  = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5" // Immortal Baron, Vorg
	gabzagulWarlordOfPainSetupSrc = "gabzagul_warlord_of_pain_test_setup"
)

func TestGabzagulWarlordOfPain(t *testing.T) {
	t.Run("holds up the turn until every able creature has attacked", func(t *testing.T) {
		scn, owner, _ := setupGabzagulWarlordOfPainTest(t)
		gabzagul := putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, owner.Player, gabzagulWarlordOfPainUID)
		ally := putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, owner.Player, gabzagulWarlordOfPainAllyUID)

		assert.Equal(t, "Gabzagul, Warlord of Pain", gabzagul.Name)
		assert.Equal(t, 5000, gabzagul.Power)
		assert.Equal(t, 6, gabzagul.ManaCost)
		assert.Equal(t, []string{civ.Darkness}, gabzagul.Civs)
		assert.True(t, gabzagul.HasFamily(family.DarkLord))

		require.NoError(t, scn.ActionEndTurn(owner))
		assert.True(t, scn.Match.IsPlayerTurn(owner.Player), "the turn cannot end while a creature can still attack")

		// Tapping stands in for those creatures having attacked.
		ally.Tapped = true
		require.NoError(t, scn.ActionEndTurn(owner))
		assert.True(t, scn.Match.IsPlayerTurn(owner.Player), "Gabzagul forces itself to attack too")

		gabzagul.Tapped = true
		require.NoError(t, scn.ActionEndTurn(owner))
		assert.False(t, scn.Match.IsPlayerTurn(owner.Player))
	})

	t.Run("forces the creatures of both players", func(t *testing.T) {
		scn, owner, opponent := setupGabzagulWarlordOfPainTest(t)
		// Gabzagul belongs to the player who is not taking the turn.
		gabzagul := putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, opponent.Player, gabzagulWarlordOfPainUID)
		ally := putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, owner.Player, gabzagulWarlordOfPainAllyUID)

		require.NoError(t, scn.ActionEndTurn(owner))
		assert.True(t, scn.Match.IsPlayerTurn(owner.Player), "an opponent's Gabzagul forces your creatures as well")
		assert.False(t, gabzagul.Tapped, "Gabzagul is not forced to attack on someone else's turn")

		ally.Tapped = true
		require.NoError(t, scn.ActionEndTurn(owner))
		assert.False(t, scn.Match.IsPlayerTurn(owner.Player))
	})

	t.Run("ignores creatures that cannot attack", func(t *testing.T) {
		scn, owner, opponent := setupGabzagulWarlordOfPainTest(t)
		putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, opponent.Player, gabzagulWarlordOfPainUID)

		// Summoned this turn, so it has summoning sickness and is not able to attack.
		owner.Player.SpawnCard(gabzagulWarlordOfPainAllyUID, match.HAND)
		for range 2 {
			owner.Player.SpawnCard(gabzagulWarlordOfPainAllyUID, match.MANAZONE)
		}
		summoned, err := scn.FindCard(owner.Player, match.HAND, gabzagulWarlordOfPainAllyUID)
		require.NoError(t, err)
		require.NoError(t, scn.ActionPlayCard(owner, summoned.ID))
		require.Equal(t, match.BATTLEZONE, summoned.Zone)

		require.NoError(t, scn.ActionEndTurn(owner))
		assert.False(t, scn.Match.IsPlayerTurn(owner.Player), "a summoning sick creature does not hold up the turn")
	})

	t.Run("stops forcing attacks once it leaves the battle zone", func(t *testing.T) {
		scn, owner, opponent := setupGabzagulWarlordOfPainTest(t)
		gabzagul := putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, opponent.Player, gabzagulWarlordOfPainUID)
		putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, owner.Player, gabzagulWarlordOfPainAllyUID)

		require.NoError(t, scn.ActionEndTurn(owner))
		require.True(t, scn.Match.IsPlayerTurn(owner.Player))

		moved, err := opponent.Player.MoveCard(gabzagul.ID, match.BATTLEZONE, match.GRAVEYARD, gabzagulWarlordOfPainSetupSrc)
		require.NoError(t, err)
		require.Equal(t, match.GRAVEYARD, moved.Zone)

		require.NoError(t, scn.ActionEndTurn(owner))
		assert.False(t, scn.Match.IsPlayerTurn(owner.Player))
	})

	t.Run("warns only once when several creatures must still attack", func(t *testing.T) {
		scn, owner, opponent := setupGabzagulWarlordOfPainTest(t)
		putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, opponent.Player, gabzagulWarlordOfPainUID)
		for range 3 {
			putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, owner.Player, gabzagulWarlordOfPainAllyUID)
		}

		warningStart, err := scn.MessageCount(owner)
		require.NoError(t, err)
		require.NoError(t, scn.ActionEndTurn(owner))
		require.True(t, scn.Match.IsPlayerTurn(owner.Player))

		// The persistent effect loop is not stopped by cancelling the context, so
		// without its own guard every able creature would warn separately.
		warnings, err := scn.Warnings(owner, warningStart)
		require.NoError(t, err)
		mustAttack := 0
		for _, warning := range warnings {
			if strings.Contains(warning, "must attack before you can end your turn") {
				mustAttack++
			}
		}
		assert.Equal(t, 1, mustAttack)
	})

	t.Run("two copies do not deadlock the turn", func(t *testing.T) {
		scn, owner, opponent := setupGabzagulWarlordOfPainTest(t)
		putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, opponent.Player, gabzagulWarlordOfPainUID)
		putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, opponent.Player, gabzagulWarlordOfPainUID)
		ally := putGabzagulWarlordOfPainTestCardInBattlezone(t, scn, owner.Player, gabzagulWarlordOfPainAllyUID)

		require.NoError(t, scn.ActionEndTurn(owner))
		require.True(t, scn.Match.IsPlayerTurn(owner.Player))

		ally.Tapped = true
		require.NoError(t, scn.ActionEndTurn(owner))
		assert.False(t, scn.Match.IsPlayerTurn(owner.Player))
	})
}

func setupGabzagulWarlordOfPainTest(t *testing.T) (*scenario.TestScenario, *match.PlayerReference, *match.PlayerReference) {
	t.Helper()

	scn := scenario.New()
	owner := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(owner.Player))

	return scn, owner, opponent
}

func putGabzagulWarlordOfPainTestCardInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	player.SpawnCard(uid, match.HAND)
	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)
	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, gabzagulWarlordOfPainSetupSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)
	return moved
}
